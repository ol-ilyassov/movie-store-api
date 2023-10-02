#!/bin/bash
set -eu

# ==================================================================================== #
# VARIABLES
# ==================================================================================== #

# List of timezones: timedatectl list-timezones.
# Set server timezone. 
TIMEZONE=Asia/Almaty
# Create user for server:
USERNAME=moviestore

# DSN formation:
read -p "Enter password for DB user: " DB_PASSWORD
# force one locale during setup to avoid extra errors:
export LC_ALL=en_US.UTF-8

# ==================================================================================== #
# SCRIPT LOGIC
# ==================================================================================== #

# Enable the "universe" repository.
add-apt-repository --yes universe
# Update all software packages.
apt update
# Set the system timezone and install all locales.
timedatectl set-timezone ${TIMEZONE}
apt --yes install locales-all
# Add the new user (and give them sudo privileges).
useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"
# Force a password to be set for the new user the first time they log in.
passwd --delete "${USERNAME}"chage --lastday 0 "${USERNAME}"
# Copy the SSH keys from the root user to the new user.
rsync --archive --chown=${USERNAME}:${USERNAME} /root/.ssh /home/${USERNAME}
# Configure the firewall to allow SSH, HTTP and HTTPS traffic.
ufw allow 22
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable
# Install fail2ban.
apt --yes install fail2ban
# Install the migrate CLI tool.
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
mv migrate.linux-amd64 /usr/local/bin/migrate
# Install PostgreSQL.
apt --yes install postgresql
# Set up the DB and create a user account with the password entered earlier.
sudo -i -u postgres psql -c "CREATE DATABASE movies_api"
sudo -i -u postgres psql -d movies_api -c "CREATE EXTENSION IF NOT EXISTS citext"
sudo -i -u postgres psql -d movies_api -c "CREATE ROLE movies_api WITH LOGIN PASSWORD '${DB_PASSWORD}'"
# Add a DSN for connecting to the database to the system-wide environment variables in the /etc/environment file.
echo "MOVIES_API_DB_DSN='postgres://movies_api:${DB_PASSWORD}@localhost/movies_api'" >> /etc/environment
# Install Caddy (see https://caddyserver.com/docs/install#debian-ubuntu-raspbian).
apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
apt update
apt --yes install caddy
# Upgrade all packages. Using the --force-confnew flag means that configuration
# files will be replaced if newer ones are available.
apt --yes -o Dpkg::Options::="--force-confnew" upgrade
echo "Script complete! Rebooting..."
reboot

# ==================================================================================== #