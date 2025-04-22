package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"starter/config"
	"starter/internal/database"
	"starter/internal/server" // Keep the import for the server
	"starter/pkg/logger"
	"github.com/pressly/goose/v3"
)

func main() {
	// Initialize logger
	logger.InitLogger()
	log := logger.GetLogger()

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
	// if err := runMigrations(config.GetConfig().DatabaseURL); err != nil {
	// 	log.Fatal().Err(err).Msg("Failed to run migrations")
	// }

	// Create and start the server, and pass the db instance
	srv := server.NewServer(db)
	log.Info().Str("addr", srv.Addr).Msg("Starting HTTP server")

	// Start the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	// Graceful shutdown on receiving SIGINT or SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	// Set a timeout for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}

// Run migrations using Goose
func runMigrations(dsn string) error {
	db, err := goose.OpenDBWithDriver("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// Run all migrations
	return goose.Up(db, "migrations")
}
