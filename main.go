package main

import (
	"context"
	"fmt"
	"github.com/perbu/csr-pattern/api"
	"github.com/perbu/csr-pattern/repo"
	"github.com/perbu/csr-pattern/service"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := run(ctx, os.Stdout); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("clean exit")
}

func run(ctx context.Context, out *os.File) error {
	// create a logger:
	logger := newLogger(out)

	// instantiate the repo:
	r, err := repo.New("repo.db", logger.With("component", "repo"))
	if err != nil {
		return fmt.Errorf("repo.New: %w", err)
	}
	defer r.Close()
	// instantiate the service:
	s, err := service.New(r, logger.With("component", "service"))
	if err != nil {
		return fmt.Errorf("service.New: %w", err)
	}
	// instantiate the API:
	a, err := api.New(s, logger.With("component", "api"))
	if err != nil {
		return fmt.Errorf("api.New: %w", err)
	}

	// start the API server:
	if err := a.Run(ctx); err != nil {
		return fmt.Errorf("api.Run: %w", err)
	}
	return nil
}

func newLogger(fh *os.File) *slog.Logger {
	lh := slog.NewTextHandler(fh, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(lh)
}
