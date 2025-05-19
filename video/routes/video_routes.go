package routes

import (
	"starter/video/controllers"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func RegisterVideoRoutes(r chi.Router, db *gorm.DB) {
	ctrl := controllers.NewVideoController(db)

	// r.With(auth.JWTAuthMiddleware).Route("/videos", func(r chi.Router) {
	r.Route("/videos", func(r chi.Router) {
		r.Post("/", ctrl.Create)
		// Handle actual file upload
		r.Post("/{id}/upload", ctrl.ProcessVideoUpload)
		r.Get("/{id}", ctrl.Get)
		r.Get("/", ctrl.GetAll)
		r.Put("/{id}", ctrl.Update)
		r.Delete("/{id}", ctrl.Delete)
		r.Get("/{id}/with-channel", ctrl.GetWithChannel)
		r.Get("/with-channel", ctrl.GetAllWithChannel)
		r.Get("/channel/{channelID}", ctrl.GetAllByChannelID)
	})
}
