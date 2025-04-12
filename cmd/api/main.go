package main

import (
	"context"
	// "embed"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"starter/config"
	"starter/internal/cache"
	"starter/internal/database"
	"starter/internal/server"
	"starter/pkg/logger"

	"github.com/pressly/goose/v3"
)

// var embedMigrations embed.FS

func main() {
	// Initialize logger
	logger.InitLogger()
	log := logger.GetLogger()
	// Initialize redis
	cache.InitRedis()

	// Load configuration
	if err := config.InitConfig(); err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Connect to database
	db, err := database.NewDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Run migrations
	if err := runMigrations(config.GetConfig().DatabaseURL); err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}

	// Create and start server
	srv := server.NewServer(db)
	log.Info().Str("addr", srv.Addr).Msg("Starting HTTP server")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}

func runMigrations(dsn string) error {
	db, err := goose.OpenDBWithDriver("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	// goose.SetBaseFS(embedMigrations)
	return goose.Up(db, "migrations")
}
