package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/dilithaw123/broccoli-backend/internal/group"
	"github.com/dilithaw123/broccoli-backend/internal/session"
	"github.com/dilithaw123/broccoli-backend/internal/user"
	"github.com/dilithaw123/broccoli-backend/internal/web"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	password, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		logger.Error("PASSWORD environment variable is required")
		os.Exit(1)
	}
	pool, err := pgxpool.New(
		context.Background(),
		"host=localhost port=5432 user=postgres password="+password+" dbname=broccoli sslmode=disable",
	)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	if err := pool.Ping(context.Background()); err != nil {
		logger.Error("Failed to ping database", "error", err)
		os.Exit(1)
	}
	logger.Info("Connected to database")
	defer pool.Close()
	userService := user.NewPgUserRepo(pool)
	groupService := group.NewPgGroupRepo(pool)
	sessionService := session.NewPgSessionRepo(pool)
	server := web.NewServer(
		pool,
		web.WithDB(pool),
		web.WithLogger(logger),
		web.WithUserService(userService),
		web.WithGroupService(groupService),
		web.WithSessionService(sessionService),
		web.WithMux(http.NewServeMux()),
	)
	if err := server.Start(":5050"); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
