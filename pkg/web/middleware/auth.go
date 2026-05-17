package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/J4yTr1n1ty/go-api-template/pkg/boot"
	"github.com/J4yTr1n1ty/go-api-template/pkg/config"
	"github.com/J4yTr1n1ty/go-api-template/pkg/reqid"
	"github.com/J4yTr1n1ty/go-api-template/pkg/web/httputil"
	"github.com/golang-jwt/jwt/v5"
)

// Scope is a named permission string used to guard routes.
type Scope string

type claimsContextKey struct{}

// Claims holds the JWT payload. Scopes controls route access.
// Embed jwt.RegisteredClaims for standard fields (sub, exp, iat, etc.).
// Extend with app-specific fields as needed.
type Claims struct {
	Scopes []string `json:"scopes"`
	jwt.RegisteredClaims
}

// RequireScope returns a Middleware that validates the JWT from the
// Authorization header and checks that the token holds all of the given
// scopes. Injects *Claims into the request context on success.
func RequireScope(required ...Scope) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := extractClaims(r)
			if err != nil {
				slog.Debug("authentication failed", "request_id", reqid.FromContext(r.Context()), "remote_addr", r.RemoteAddr)
				httputil.ErrorJSON(w, r, http.StatusUnauthorized, "unauthorized")
				return
			}

			for _, s := range required {
				if !slices.Contains(claims.Scopes, string(s)) {
					slog.Debug("scope check failed", "request_id", reqid.FromContext(r.Context()), "remote_addr", r.RemoteAddr, "required", s, "scopes", claims.Scopes)
					httputil.ErrorJSON(w, r, http.StatusForbidden, "forbidden")
					return
				}
			}

			ctx := context.WithValue(r.Context(), claimsContextKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ClaimsFromContext retrieves the injected Claims from a request context.
// Returns nil if no claims are present (e.g. on public routes).
func ClaimsFromContext(ctx context.Context) *Claims {
	v, _ := ctx.Value(claimsContextKey{}).(*Claims)
	return v
}

func extractClaims(r *http.Request) (*Claims, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return nil, errors.New("missing Authorization header")
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, errors.New("malformed Authorization header")
	}

	secret := boot.Environment.GetEnv(config.EnvJwtSecret)
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(parts[1], claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
