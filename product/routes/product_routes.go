package routes

import (
	"starter/auth"
	"starter/product/controllers"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func RegisterProductRoutes(r chi.Router, db *gorm.DB) {
	ctrl := controllers.NewProductController(db)

	r.With(auth.JWTAuthMiddleware).Route("/products", func(r chi.Router) {
		r.Post("/", ctrl.Create)
		r.Get("/{id}", ctrl.Get)
		r.Put("/{id}", ctrl.Update)
		r.Delete("/{id}", ctrl.Delete)
		r.Post("/{id}/image", ctrl.UploadImage)
	})
}
