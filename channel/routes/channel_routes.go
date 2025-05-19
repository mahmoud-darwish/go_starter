package routes

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"starter/channel/controllers"
	"starter/auth"
	

)

func RegisterChannelRoutes(r *chi.Mux, db *gorm.DB) {
	ctrl := controllers.NewChannelController(db)

	r.With(auth.JWTAuthMiddleware).Route("/channels", func(r chi.Router) {
		r.Post("/", ctrl.Create)
		r.Get("/", ctrl.FindAll)  
		r.Get("/{id}", ctrl. FindByID)
		r.Put("/{id}", ctrl.Update)
		r.Delete("/{id}", ctrl.Delete)
		r.Post("/{id}/image", ctrl.UploadImage)
	})
}
