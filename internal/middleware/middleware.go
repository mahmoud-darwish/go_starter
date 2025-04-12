package middleware

import (
	"net/http"
	"time"

	"starter/pkg/logger"

	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

// logger middleware
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLogger()
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Msg("Received HTTP request")
		next.ServeHTTP(w, r)
	})
}

func SetupMiddleware(handler http.Handler) http.Handler {
	log := logger.GetLogger()

	// CORS
	corsMiddleware := cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"http://localhost:3000", "https://yourfrontend.com"},
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	log.Info().Msg("CORS middleware configured")

	// rate limit 100 requests per minute per ip
	rateLimitMiddleware := httprate.LimitByIP(100, 1*time.Minute)
	log.Info().Msg("Rate limiting middleware configured")
	// for user-specific limits, use httprate.LimitByRealIP with JWT claims

	// 	// CSRF protection in future
	// 	csrfSecret := os.Getenv("CSRF_SECRET")
	// 	if csrfSecret == "" {
	// 		log.Fatal().Msg("CSRF_SECRET not set")
	// 	}
	// 	csrfMiddleware := csrf.Protect(
	// 		[]byte(csrfSecret),
	// 		csrf.Secure(false),
	// 		csrf.Path("/"),
	// 	)
	// 	log.Info().Msg("CSRF middleware configured")

	// Apply middleware
	handler = corsMiddleware(handler)
	handler = rateLimitMiddleware(handler)
	// 	handler = csrfMiddleware(handler)

	return handler
}
