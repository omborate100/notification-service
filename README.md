# Notification Service

A production-ready notification microservice built with **Go** that
provides template-based email delivery for applications.

## Features

-   Template-based HTML emails
-   Dynamic variable rendering
-   SMTP (Local) and Brevo API (Production)
-   Environment-based provider selection
-   PostgreSQL template & tracking storage
-   Reply-To support
-   Render deployment ready

## Architecture

``` text
Client Application
        |
        v
Notification Service
  ├── Handler
  ├── Service
  ├── Repository
  └── Mail Factory
        |
   +----+----+
   |         |
 SMTP     Brevo API
(Local)  (Production)
```

## Project Structure

``` text
notification-service/
├── cmd/
├── config/
├── internal/
│   ├── handler/
│   ├── model/
│   ├── renderer/
│   ├── repository/
│   ├── routes/
│   ├── service/
│   └── mail/
│       ├── sender.go
│       ├── factory.go
│       ├── smtp/
│       └── brevo/
├── migrations/
├── templates/
└── README.md
```

## Environment

### Local

``` env
APP_ENV=LOCAL
SMTP_HOST=smtp-relay.brevo.com
SMTP_PORT=587
SMTP_USERNAME=...
SMTP_PASSWORD=...

FROM_EMAIL=BandhanBio <noreply@omborate.in>
REPLY_TO_EMAIL=BandhanBio Support <admin.bandhanbio@gmail.com>
```

### Production

``` env
APP_ENV=PROD
BREVO_API_KEY=your_api_key

FROM_EMAIL=BandhanBio <noreply@omborate.in>
REPLY_TO_EMAIL=BandhanBio Support <admin.bandhanbio@gmail.com>
```

## Email Flow

``` text
Application
    |
REST API
    |
Template Rendering
    |
Mail Factory
    |
+---+---------+
|             |
SMTP      Brevo API
```

## Running Locally

``` bash
go mod tidy
go run cmd/server/main.go
```

## Deployment

Deploy on Render.

-   Local uses SMTP.
-   Production uses the Brevo REST API to avoid blocked SMTP ports on
    Render Free.

## Tech Stack

-   Go
-   PostgreSQL
-   Brevo
-   SMTP
-   Render

## Future Improvements

-   SMS
-   Push Notifications
-   Bulk Emails
-   Retry mechanism
-   Webhooks
-   Metrics
-   Swagger

## License

MIT
