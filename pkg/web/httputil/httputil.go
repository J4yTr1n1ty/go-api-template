package httputil

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/J4yTr1n1ty/go-api-template/pkg/config"
	"github.com/J4yTr1n1ty/go-api-template/pkg/reqid"
	"github.com/jackc/pgx/v5/pgconn"
)

// InternalError sends a 500 response, logging the error.
// Debug mode: full error string in response body.
// Production mode: generic "Internal Server Error" — no internal details exposed.
func InternalError(w http.ResponseWriter, r *http.Request, err error) {
	rid := reqid.FromContext(r.Context())
	slog.Error("internal server error", "request_id", rid, "err", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	msg := http.StatusText(http.StatusInternalServerError)
	if config.IsDebug() {
		msg = err.Error()
	}
	json.NewEncoder(w).Encode(struct { //nolint:errcheck
		Error     string `json:"error"`
		RequestID string `json:"request_id,omitempty"`
	}{msg, rid})
}

// ErrorJSON writes a JSON error response with the given status code and message.
// request_id from context is included in the response body when present.
func ErrorJSON(w http.ResponseWriter, r *http.Request, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(struct { //nolint:errcheck
		Error     string `json:"error"`
		RequestID string `json:"request_id,omitempty"`
	}{msg, reqid.FromContext(r.Context())})
}

// ValidationErrorJSON writes a 422 JSON response with field-level validation errors.
// fields is a map of field name to error message.
func ValidationErrorJSON(w http.ResponseWriter, r *http.Request, fields map[string]string) {
	rid := reqid.FromContext(r.Context())
	slog.Warn("validation failed",
		"request_id", rid,
		"method", r.Method,
		"path", r.URL.Path,
		"fields", fields,
	)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(struct { //nolint:errcheck
		Error     string            `json:"error"`
		Fields    map[string]string `json:"fields"`
		RequestID string            `json:"request_id,omitempty"`
	}{"validation failed", fields, rid})
}

// DBError inspects err for known PostgreSQL constraint violations and writes an
// appropriate JSON response. Returns true if the error was handled, false otherwise.
func DBError(w http.ResponseWriter, r *http.Request, err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	switch pgErr.Code {
	case "23505": // unique_violation
		ErrorJSON(w, r, http.StatusConflict, "already exists")
		return true
	case "23503": // foreign_key_violation
		ErrorJSON(w, r, http.StatusBadRequest, "referenced resource not found")
		return true
	}
	return false
}

// PathID parses a uint path parameter by name. On failure writes a 400 JSON response and returns false.
func PathID(w http.ResponseWriter, r *http.Request, param string) (uint, bool) {
	v := r.PathValue(param)
	if v == "" {
		ErrorJSON(w, r, http.StatusBadRequest, "missing "+param+" parameter")
		return 0, false
	}
	n, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		ErrorJSON(w, r, http.StatusBadRequest, "invalid "+param+" parameter")
		return 0, false
	}
	return uint(n), true
}

// PaginationParams parses skip and take query parameters, applying the given defaults and cap.
// On failure writes a 400 JSON response and returns false.
func PaginationParams(w http.ResponseWriter, r *http.Request, defaultTake, maxTake int) (skip, take int, ok bool) {
	qp := r.URL.Query()
	skip = 0
	take = defaultTake

	if v := qp.Get("skip"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			ErrorJSON(w, r, http.StatusBadRequest, "invalid skip")
			return 0, 0, false
		}
		skip = n
	}

	if v := qp.Get("take"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			ErrorJSON(w, r, http.StatusBadRequest, "invalid take")
			return 0, 0, false
		}
		take = n
	}

	if take > maxTake {
		take = maxTake
	}

	return skip, take, true
}
