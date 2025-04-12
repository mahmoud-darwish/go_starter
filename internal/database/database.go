package database

import (
	"starter/config"
	"starter/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
	cfg := config.GetConfig()
	log := logger.GetLogger()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		return nil, err
	}

	log.Info().Msg("Database connection established")
	return db, nil
}
