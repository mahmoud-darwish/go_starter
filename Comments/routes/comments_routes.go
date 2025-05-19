package routes

import (
	"starter/auth"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"starter/Comments/controllers"
	
	

)

func RegisterCommentsRoutes(r *chi.Mux, db *gorm.DB) {
	ctrl := controllers.NewCommentsController(db)

	r.With(auth.JWTAuthMiddleware).Route("/comments", func(r chi.Router) {
		r.Post("/", ctrl.Create)
		r.Get("/", ctrl.FindAll)  
		r.Get("/{id}", ctrl. FindByID)
	})
}
