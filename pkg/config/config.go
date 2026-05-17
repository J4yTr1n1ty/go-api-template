package config

import (
	"github.com/J4yTr1n1ty/go-api-template/pkg/boot"
)

const (
	EnvDebug        = "DEBUG"
	EnvPort         = "PORT"
	EnvPostgresHost = "POSTGRES_HOST"
	EnvPostgresPort = "POSTGRES_PORT"
	EnvPostgresUser = "POSTGRES_USER"
	EnvPostgresPass = "POSTGRES_PASS"
	EnvPostgresDB   = "POSTGRES_DB"
	EnvLogFormat    = "LOG_FORMAT"  // "json" (default) or "text"
	EnvJwtSecret    = "JWT_SECRET"  // required when using RequireScope middleware
)

func IsDebug() bool {
	return boot.Environment.GetEnvBool(EnvDebug)
}
