# Go API Template

A clean, minimal Go REST API template with PostgreSQL, Docker support, and structured middleware.

## Features

- **PostgreSQL** integration with GORM
- **Docker** support with multi-stage builds
- **Middleware** stack (logging, recovery)
- **Environment-based** configuration
- **Clean architecture** with separated concerns

## Quick Start

### Using Docker (Recommended)

```bash
docker-compose up --build
```

### Local Development

```bash
# Copy environment file
cp .env.example .env

# Install dependencies
go mod tidy

# Run the application
go run cmd/server/main.go
```

## Environment Variables

Create a `.env` file with:

```bash
PORT=8080
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASS=postgres
POSTGRES_DB=rest_api
DEBUG=false
```

## Project Structure

```
├── cmd/server/          # Application entry point
├── pkg/
│   ├── boot/           # Bootstrapping (env, database)
│   ├── config/         # Configuration constants
│   ├── models/         # Data models
│   └── web/
│       ├── handlers/   # HTTP handlers
│       ├── middleware/ # HTTP middleware
│       └── router.go   # Route definitions
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

## API Endpoints

- `GET /example` - Example endpoint returning data from database

## Adding New Features

### 1. Add a Model

```go
// pkg/models/user.go
type User struct {
    gorm.Model
    Name  string
    Email string
}
```

### 2. Create a Handler

```go
// pkg/web/handlers/user/user.go
func (h *Handler) GetUsers() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Handler logic
    }
}
```

### 3. Register Routes

```go
// pkg/web/router.go
userHandler := user.NewHandler()
mux.Handle("GET /users", userHandler.GetUsers())
```

### 4. Run Migrations

Add your model to the AutoMigrate call in `cmd/server/main.go`:

```go
err := boot.DB.AutoMigrate(&models.User{})
```

## Development

```bash
# Format code
go fmt ./...

# Run tests
go test ./...

# Build binary
go build -o bin/server cmd/server/main.go
```

## License

This project is licensed under the AGPL v3 License.
