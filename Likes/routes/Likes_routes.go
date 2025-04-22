package routes

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"starter/Likes/controllers"
	
	

)

func RegisterLikesRoutes(r *chi.Mux, db *gorm.DB) {
	ctrl := controllers.NewLikeController(db)

	r.Route("/likes", func(r chi.Router) {
		r.Post("/", ctrl.Create)
		r.Get("/", ctrl.FindAll)  
		r.Get("/{id}", ctrl.FindByID)
		r.Delete("/{id}", ctrl.Delete)
	})
}
