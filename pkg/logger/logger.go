package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func InitLogger() {
	logger = zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger().
		Level(zerolog.DebugLevel)
}

func GetLogger() *zerolog.Logger {
	return &logger
}
