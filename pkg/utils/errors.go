package utils

import (
	"encoding/json"
	"net/http"

	"starter/pkg/logger"
)

func RespondWithError(w http.ResponseWriter, status int, message string) {
	log := logger.GetLogger()
	log.Error().Int("status", status).Str("error", message).Msg("HTTP error response")

	response := map[string]string{"error": message}
	RespondWithJSON(w, status, response)
}

func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log := logger.GetLogger()
		log.Error().Err(err).Msg("Failed to encode JSON response")
	}
}
