# Go API Template

Clean Go REST API template with PostgreSQL, JWT auth, structured logging, and Swagger docs.

## Features

- **PostgreSQL** via GORM
- **Structured logging** with `log/slog` (JSON default, text via `LOG_FORMAT=text`)
- **JWT auth** middleware with scope-based route guards
- **Request ID** propagated through context, logs, and error responses
- **Validation** via go-playground/validator with JSON field name errors
- **Swagger UI** at `/swagger/` — regenerate with `go generate ./cmd/server/...`
- **Pagination** — `PaginationResult[T]` generic + `PaginationMeta` in `pkg/dtos`
- **httputil** — `InternalError`, `ErrorJSON`, `ValidationErrorJSON`, `DBError`, `PathID`, `PaginationParams`
- **Docker** multi-stage build, non-root user
- **Middleware** stack: RequestID → Logging → Recovery

## Quick Start

```bash
cp .env.example .env
# edit .env

docker-compose up --build
# or
go run cmd/server/main.go
```

## Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `PORT` | yes | — | HTTP listen port |
| `POSTGRES_HOST` | yes | — | |
| `POSTGRES_PORT` | yes | — | |
| `POSTGRES_USER` | yes | — | |
| `POSTGRES_PASS` | yes | — | |
| `POSTGRES_DB` | yes | — | |
| `DEBUG` | no | `false` | Debug logging + full error bodies |
| `LOG_FORMAT` | no | `json` | `json` or `text` |
| `JWT_SECRET` | when using auth | — | HMAC secret for JWT signing |

## Project Structure

```
cmd/server/          entry point, slog setup, swagger annotations
docs/                generated swagger spec (go generate ./cmd/server/...)
pkg/
  boot/              env loading, DB connection
  config/            env var name constants, IsDebug()
  dtos/              request/response structs, PaginationResult[T]
  models/            GORM models
  reqid/             request ID context storage
  validate/          validate.Struct wrapper, FieldErrors type
  web/
    handlers/        HTTP handlers
    httputil/        error response helpers
    middleware/      RequestID, Logging, Recovery, RequireScope
    router.go        route registration, scope()/wrap() helpers
```

## Adding a Resource

### 1. Model

```go
// pkg/models/thing.go
type Thing struct {
    gorm.Model
    Name string
}
```

### 2. DTOs

```go
// pkg/dtos/thing.go
type CreateThingRequest struct {
    Name string `json:"name" validate:"required,max=255"`
}

type ThingResponse struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
}

// concrete type required for swag (cannot use generics directly)
type ThingListResponse struct {
    Data     []ThingResponse `json:"data"`
    Metadata PaginationMeta  `json:"metadata"`
}
```

### 3. Handler

```go
// pkg/web/handlers/thing/handler.go
func (h *Handler) Create() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req dtos.CreateThingRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            httputil.ErrorJSON(w, r, http.StatusBadRequest, "invalid request body")
            return
        }
        if err := validate.Struct(req); err != nil {
            httputil.ValidationErrorJSON(w, r, err.(validate.FieldErrors))
            return
        }
        // ...
    }
}
```

### 4. Scopes + Routes

```go
// pkg/web/router.go
const ScopeThingWrite middleware.Scope = "thing:write"

mux.Handle("POST /things", wrap(scope(ScopeThingWrite), thingHandler.Create))
```

### 5. AutoMigrate

```go
// cmd/server/main.go
boot.DB.AutoMigrate(&models.Thing{})
```

## Development

```bash
go fmt ./...
go vet ./...
go build ./...
go generate ./cmd/server/...   # regenerate swagger docs
```

## License

AGPL v3.
