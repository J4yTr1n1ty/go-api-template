package config

import (
	"github.com/J4yTr1n1ty/go-api-template/pkg/boot"
)

const (
	EnvPort         = "PORT"
	EnvPostgresHost = "POSTGRES_HOST"
	EnvPostgresPort = "POSTGRES_PORT"
	EnvPostgresUser = "POSTGRES_USER"
	EnvPostgresPass = "POSTGRES_PASS"
	EnvPostgresDB   = "POSTGRES_DB"
)

func IsDebug() bool {
	return boot.Environment.GetEnvBool("DEBUG")
}
