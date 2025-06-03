package main

import (
	"cmp"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcusolsson/cookhub/api"
	dbpkg "github.com/marcusolsson/cookhub/db"
	"github.com/marcusolsson/cookhub/ui"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var (
		port = cmp.Or(os.Getenv("PORT"), "8080")
		dsn  = os.Getenv("DATABASE_URL")
	)

	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger.Info("Attempting to connect to database...")

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.Error("Failed to connect to the database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Error("Failed to ping the database", "error", err)
		os.Exit(1)
	}

	logger.Info("Successfully connected to the database")

	store := dbpkg.NewStore(pool)

	if err := store.CreateJobsTable(ctx); err != nil {
		logger.Error("Failed to create jobs table", "error", err)
		os.Exit(1)
	}

	if err := store.CreateFilesTable(ctx); err != nil {
		logger.Error("Failed to create files table", "error", err)
		os.Exit(1)
	}

	apisrv := api.NewRouter(store)
	uisrv := ui.NewRouter(store)

	r := chi.NewRouter()
	r.Mount("/api", apisrv)
	r.Mount("/ui", uisrv)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		logger.Info("Starting server...", "addr", srv.Addr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", "error", err)
		}
	}()

	waitForShutdownSignal()

	logger.Info("Received shutdown signal, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Failed to shutdown server", "err", err)
		os.Exit(1)
	}
}

func waitForShutdownSignal() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
