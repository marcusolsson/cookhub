package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/marcusolsson/cookhub/db/sqlc"
	"github.com/patrickmn/go-cache"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var (
		cfg    = loadConfigFromEnv()
		logger = newLogger(cfg)
		ctx    = context.Background()
	)

	logger.Info("Connecting to the database...")

	pool, err := pgxpool.New(ctx, cfg.DB.URL)
	if err != nil {
		exitWithError(logger, "Failed to connect to the database", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		exitWithError(logger, "Failed to ping the database", err)
	}

	logger.Info("Successfully connected to the database")

	s := &Server{
		pool:     pool,
		logger:   logger,
		ghClient: newGitHubClient(cfg.GitHub.Token),
		db:       db.New(pool),
		c:        cache.New(cache.NoExpiration, cache.NoExpiration),
	}

	srv := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: s.Router(),
	}

	go func() {
		logger.Info("Starting server...", "addr", srv.Addr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			exitWithError(logger, "Failed to start server", err)
		}
	}()

	waitForShutdownSignal()

	logger.Info("Received shutdown signal, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		exitWithError(logger, "Failed to gracefully shutdown server", err)
	}
}

// waitForShutdownSignal waits for an interrupt or termination signal to gracefully shut down the server.
func waitForShutdownSignal() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}

// exitWithError logs the error message and exits the program with a non-zero status code.
func exitWithError(logger *slog.Logger, message string, err error) {
	logger.Error(message, "error", err)
	os.Exit(1)
}
