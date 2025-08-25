# goconda — Beego v2.3.8 boilerplate (Go 1.25)

A reusable, extendible boilerplate for building web apps with Beego v2.3.8 using Go 1.25.0.

## Features

- Authentication: **Register** and **Login** with email, firstname, lastname, password (User primary key = **email**)
- Token-based auth (JWT) **and** session-based auth
- Middleware example for protecting routes
- Controllers with CRUD (sample `Item` resource) + protected `GET /api/v1/users/me`
- ORM + DB integration (Beego ORM)
- Routing via `web.NewNamespace`
- Static files served from `/static`
- File upload endpoint (`POST /api/v1/upload`)
- **Dev** (SQLite) and **Prod** (MariaDB/MySQL) configs
- Unit tests for JWT, hashing, validators

## Quick start

```bash
cd goconda

export APP_ENV=dev   # or prod

go mod tidy
go run .
```

Server listens on `:8080` by default (see `conf/app.*.conf`).

## API

### Auth
- `POST /api/v1/auth/register` JSON:
  ```json
  {
    "email": "you@example.com",
    "firstname": "Ada",
    "lastname": "Lovelace",
    "password": "secret123"
  }
  ```
- `POST /api/v1/auth/login` JSON:
  ```json
  {
    "email": "you@example.com",
    "password": "secret123"
  }
  ```
  Returns `{ token, user }` and also sets a session cookie (`email`).

Include `Authorization: Bearer <token>` **or** use the session cookie for authenticated routes.

### Users
- `GET /api/v1/users/me` — returns your user profile (auth required).

### Items (CRUD, auth required)
- `GET /api/v1/items?limit=20&offset=0`
- `POST /api/v1/items` JSON: `{ "name": "...", "description": "..." }`
- `GET /api/v1/items/:id`
- `PUT /api/v1/items/:id`
- `DELETE /api/v1/items/:id`

### Uploads
- `POST /api/v1/upload` form-data field `file`

## Configuration

- `APP_ENV=dev` uses `conf/app.dev.conf` with **SQLite**
- `APP_ENV=prod` uses `conf/app.prod.conf` with **MariaDB/MySQL**

Edit JWT secret, DB DSN, and upload dir in the respective config file.

## Notes

- SQLite driver requires CGO. For MariaDB, ensure DSN is valid and DB exists.
- The `Item` model uses `owner_email` as FK referencing `users.email`.

## Testing

```bash
go test ./...
```

---


<!-- docker run -p 8080:8080 \    
  --add-host=host.docker.internal:host-gateway \
  -e BEEGO_APP_CONFIG_PATH=/app/conf/app.prod.conf \
  -e APP_ENV=prod \
  -e DB_DSN=$DB_DSN \
  -e JWT_SECRET="supersecret" \
  goconda -->

## Email Verification

New endpoints:
- `POST /api/v1/auth/send-verification` with form field `email`
- `GET /api/v1/auth/verify?token=...`

The verification email is simulated by writing the token to the HTTP response (replace with SMTP in production).

Models added:
- `EmailVerificationToken`
- `VerifiedUser`

Run tests:
```bash
go test ./...
```


## Email Sending

Configure SMTP in `conf/*.conf`:
```
[smtp]
host = smtp.example.com
port = 587
username = no-reply@example.com
password = yourpassword
from = "No Reply <no-reply@example.com>"
default_subject = "Notification"
```

Use in code:
```go
import "github.com/mymi14s/goconda/utils/mailer"

_ = mailer.SendEmail("<b>Hello</b> world", []string{"user@example.com"})
```

## Task Scheduler

A thin wrapper over `robfig/cron` with a global registry.

```go
import "github.com/mymi14s/goconda/utils/scheduler"

scheduler.Start()
scheduler.Register("say-hello", "*/10 * * * * *", func(){ fmt.Println("hello every 10s") })
```

## Roles & Permissions

Models: `Role`, `UserRole`, `Permission`. Check within controllers:
```go
if !c.RequirePermission("items", "read") { return }
```
Superusers bypass checks automatically.

## Initial Administrator

Set in config to create a superuser on startup:
```
[admin]
email = admin@example.com
password = changeme
```


### Account management endpoints
- `POST /api/v1/auth/forgot-password` (email)
- `POST /api/v1/auth/reset-password` (token, new_password)
- `POST /api/v1/auth/change-password` (current_password, new_password) — requires auth
- `POST /api/v1/auth/change-email` (password, new_email) — requires auth, superuser can bypass password check

Superusers bypass email verification automatically. The initial admin password is set **only once** on first creation; later you can change it via `change-password`.
