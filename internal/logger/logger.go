package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/marcusolsson/cookhub/internal/config"
)

func New(cfg config.Config) *slog.Logger {
	if cfg.Env == "development" {
		return slog.New(tint.NewHandler(os.Stdout, nil))
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
