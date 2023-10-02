
# Movie Store API

Server side application on movie store functionality with defined HTTP API.

## Technological stack

- Built in Go version 1.21.1
- Uses [httprouter](https://github.com/julienschmidt/httprouter) a high performance HTTP request router
- Uses [lib/pq](https://github.com/lib/pq) a pure go postgres driver for database/sql
- Uses [migrate](https://github.com/golang-migrate/migrate) a cli tool and go library for database migrations 
- Uses [go-mail/mail](https://github.com/go-mail/mail) a simple and efficient package to send emails and based on fork of gomail

## Instruments list:

- Postgresql - a high-performance SQL database
- Mailtrap - a simple SMTP server platform
- Caddy - reverse-proxy
- hey - cmd tool for load tests on API requests.