package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/J4yTr1n1ty/go-api-template/pkg/boot"
	"github.com/J4yTr1n1ty/go-api-template/pkg/config"
	"github.com/J4yTr1n1ty/go-api-template/pkg/models"
	"github.com/J4yTr1n1ty/go-api-template/pkg/web"
	"github.com/J4yTr1n1ty/go-api-template/pkg/web/middleware"
)

func main() {
	boot.LoadEnv()

	level := slog.LevelInfo
	if config.IsDebug() {
		level = slog.LevelDebug
	}
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				if src, ok := a.Value.Any().(*slog.Source); ok {
					a.Value = slog.StringValue(fmt.Sprintf("%s:%d", filepath.Base(src.File), src.Line))
				}
			}
			return a
		},
	}

	var handler slog.Handler
	if os.Getenv(config.EnvLogFormat) == "text" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	slog.SetDefault(slog.New(handler))

	slog.Info("starting server")
	boot.ConnectDB()

	if err := boot.DB.AutoMigrate(&models.Example{}); err != nil {
		slog.Error("auto migration failed", "err", err)
		os.Exit(1)
	}

	router := web.SetupRoutes()

	stack := middleware.CreateStack(middleware.RequestID, middleware.Logging, middleware.Recovery)

	server := http.Server{
		Addr:    ":" + boot.Environment.GetEnv(config.EnvPort),
		Handler: stack(router),
	}

	slog.Info("listening", "port", boot.Environment.GetEnv(config.EnvPort))
	if err := server.ListenAndServe(); err != nil {
		slog.Error("server failed", "err", err)
		os.Exit(1)
	}
}
