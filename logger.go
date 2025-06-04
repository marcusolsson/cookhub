package main

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func New(cfg Config) *slog.Logger {
	if cfg.Env == "development" {
		return slog.New(tint.NewHandler(os.Stdout, nil))
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
