
# Movie Store API

Server side application on movie store functionality with defined HTTP API.

## Technological stack

- Built in Go version 1.21.1
- Uses [httprouter](https://github.com/julienschmidt/httprouter) a high performance HTTP request router
- Uses [lib/pq](https://github.com/lib/pq) a pure go postgres driver for database/sql
- Uses [go-mail/mail](https://github.com/go-mail/mail) a simple and efficient package to send emails and based on fork of gomail
- Uses [tomasen/realip](https://github.com/tomasen/realip) a golang library that help's to get client's real public ip adress from http request headers

## Instruments list:

- [PostgreSQL](https://www.postgresql.org/) - a high-performance SQL database
- [migrate](https://github.com/golang-migrate/migrate) a cli tool and go library for database migrations
- [Mailtrap](https://mailtrap.io/) - a simple SMTP server platform
- [hey](https://github.com/rakyll/hey) - cmd tool for load tests on API requests
- [Caddy](https://caddyserver.com/docs/) - a powerful, extensible web server or proxy