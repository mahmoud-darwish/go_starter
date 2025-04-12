package routes

import (
	"starter/auth"
	"starter/internal/middleware"
	"starter/user/controllers"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func RegisterUserRoutes(r chi.Router, db *gorm.DB) {
	userCtrl := controllers.NewUserController(db)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Post("/users/register", userCtrl.Register)
		r.Post("/users/login", userCtrl.Login)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(auth.JWTAuthMiddleware)
		r.Get("/users/{id}", userCtrl.GetUser)
		r.Put("/users/{id}", userCtrl.UpdateUser)
		r.Delete("/users/{id}", userCtrl.DeleteUser)
	})
}
