package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/matthosch/go_ecommerce_api/internal/env"
)

func main() {
	ctx := context.Background()

	// Configuration
	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "user=postgres password=postgres host=localhost dbname=go_ecommerce_api sslmode=disable"),
		},
	}

	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Database
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}

	defer conn.Close(ctx)

	logger.Info("connected to database", "dsn", cfg.db.dsn)

	// Application
	api := application{
		config: cfg,
		db:     conn,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("error starting server", "error", err)
		os.Exit(1)
	}
}
