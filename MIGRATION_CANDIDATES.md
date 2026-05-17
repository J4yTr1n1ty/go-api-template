# Migration Candidates from match-api

Things worth porting to this template, ordered by value.
Source repo: `../match-api`

---

## Priority 1 — foundational, do these first

### 1. `slog` structured logging

Replace `log` package with `log/slog`. JSON output by default, text fallback via `LOG_FORMAT` env.
Includes multi-handler fan-out for optional log shipping (e.g. Seq).

| What | Source |
|------|--------|
| slog setup, handler selection, LOG_FORMAT env | `../match-api/cmd/server/main.go` (lines ~40–70) |
| Optional Seq HTTP/CLEF handler | `../match-api/pkg/seqhandler/seqhandler.go` |
| LOG_FORMAT config constant | `../match-api/pkg/config/config.go` |

### 2. Request ID middleware + `reqid` package

Crypto-random 8-byte hex ID per request. Propagated via `X-Request-ID` header.
Stored in context so handlers and downstream logs can reference it.

| What | Source |
|------|--------|
| Context storage helpers | `../match-api/pkg/reqid/reqid.go` |
| Middleware that generates/propagates ID | `../match-api/pkg/web/middleware/request_id.go` |

Update `middleware.go` to include it in `CreateStack`.

### 3. `httputil` package

Structured HTTP error responses. Every API needs this; painful to retrofit later.

| What | Source |
|------|--------|
| Full httputil package | `../match-api/pkg/web/httputil/httputil.go` |

Key functions to port (adapt or strip app-specific ones):
- `InternalError(w, r, err)` — debug-aware 500
- `ErrorJSON(w, status, msg)` — generic `{error: string}`
- `ValidationErrorJSON(w, r, fields)` — 422 with field errors
- `DBError(w, err) bool` — PG constraint → HTTP status
- `PathID(w, r, param) (uint, bool)` — parse uint path param
- `PaginationParams(w, r, defaultTake, maxTake)` — skip/take with caps

### 4. Health check endpoint

One line. Should be in every template.

| What | Source |
|------|--------|
| Route + inline handler | `../match-api/pkg/web/router.go` line 57 |

Add `GET /health` returning `200 OK` to `SetupRoutes()`.

### 5. `DEBUG` env var + debug-aware recovery

`DEBUG=true` → stack traces in error responses, verbose logging.
`DEBUG=false` (default) → generic errors only.

| What | Source |
|------|--------|
| `IsDebug()` helper, `EnvDebug` constant | `../match-api/pkg/config/config.go` |
| Debug mode in recovery middleware | `../match-api/pkg/web/middleware/recover.go` |
| Debug mode in `InternalError` | `../match-api/pkg/web/httputil/httputil.go` |
| `DEBUG` in required env list | `../match-api/pkg/boot/environment.go` |

---

## Priority 2 — common needs

### 6. JWT auth middleware (`RequireScope`)

Scope-based bearer token auth. Not every project uses it, but it's the most common
pattern and hard to retrofit cleanly once routes are established.

| What | Source |
|------|--------|
| JWT middleware | `../match-api/pkg/web/middleware/auth.go` |
| `scope()` + `wrap()` route helpers | `../match-api/pkg/web/router.go` lines 22–28 |

Dependencies: `github.com/golang-jwt/jwt/v5`, `JWT_SECRET` env var.

### 7. `go-playground/validator` integration

Field-level validation with structured JSON error output.
Pairs with `ValidationErrorJSON` from item 3.

| What | Source |
|------|--------|
| Validator wrapper | `../match-api/pkg/validate/validate.go` |

Dependencies: `github.com/go-playground/validator/v10` + locale packages.

### 8. DTOs pattern (`pkg/dtos/`)

Separate request/response structs from GORM models.
Prevents leaking DB internals and makes validation tags cleaner.

| What | Source |
|------|--------|
| Example DTO files | `../match-api/pkg/dtos/` (pick one simple one as a template) |

No direct file to copy — just establish the pattern with one example DTO.

### 9. `cmd/seed/` — db bootstrap command pattern

Shows how to structure a second `cmd/` entry point.
Content is app-specific but the pattern (env load → db connect → seed) is reusable.

| What | Source |
|------|--------|
| Seed command | `../match-api/cmd/seed/main.go` |

### 10. Swagger/swaggo integration

`swag init` + handler annotations + `/swagger/` route.
Boilerplate-heavy to add later; better to scaffold it from day one.

| What | Source |
|------|--------|
| Handler comment annotations (examples) | `../match-api/pkg/web/handlers/` (any handler file) |
| Swagger route setup | `../match-api/pkg/web/router.go` (swagger block near top) |
| `go generate` hook | `../match-api/cmd/server/main.go` (top-of-file `//go:generate`) |

Dependencies: `github.com/swaggo/swag`, `github.com/swaggo/http-swagger`.

---

## Priority 3 — nice to have

### 11. GitHub Actions CI

Build + test + optional swag doc regen.

| What | Source |
|------|--------|
| Build/test workflow | `../match-api/.github/workflows/go.yml` |
| Release workflow (GoReleaser) | `../match-api/.github/workflows/release.yml` |

### 12. GoReleaser + pre-built binary Dockerfile

Current template builds inside Docker. match-api pattern: build outside with goreleaser,
COPY pre-built binary into minimal alpine image.

| What | Source |
|------|--------|
| GoReleaser config | `../match-api/.goreleaser.yaml` |
| Production Dockerfile | `../match-api/Dockerfile` |

### 13. `scope()` + `wrap()` route helpers

Cleaner than wrapping each handler inline once you have auth middleware.
Only useful after item 6 (JWT auth) is ported.

| What | Source |
|------|--------|
| Helper functions | `../match-api/pkg/web/router.go` lines 22–28 |

### 14. Seq / MultiHandler log fan-out

Niche (Seq-specific) but the pattern of fanning out to multiple slog handlers is generic.
Worth porting the MultiHandler scaffolding even if Seq integration is dropped.

| What | Source |
|------|--------|
| MultiHandler + batch CLEF shipper | `../match-api/pkg/seqhandler/seqhandler.go` |

### 15. Go version bump

Template is on 1.24.3; match-api is on 1.26.

| What | Source |
|------|--------|
| `go.mod` toolchain line | `../match-api/go.mod` |
| Builder image in Dockerfile | current `Dockerfile` |
