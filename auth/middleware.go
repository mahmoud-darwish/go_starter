package auth

import (
	"net/http"
	"context"
	"starter/pkg/utils"
)
type contextKey string

const UserIDKey contextKey = "user_id"
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		} else {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid Authorization header")
			return
		}
		token, claims, err := VerifyJWT(tokenString)

		if err != nil || !token.Valid {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}
		userIDFloat, ok := claims["sub"].(float64) // jwt claims always come as float64
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token payload")
			return
		}
		userID := uint(userIDFloat)
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
