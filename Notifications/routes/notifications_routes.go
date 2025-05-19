package routes

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"starter/Notifications/controllers"
	"starter/auth"
	
	

)

func RegisterNotificationRoutes(r *chi.Mux, db *gorm.DB) {
	ctrl := controllers.NewNotificationController(db)

	r.With(auth.JWTAuthMiddleware).Route("/notifications", func(r chi.Router) {
		r.Post("/", ctrl.Create)
		r.Get("/", ctrl.FindAll)  
		r.Get("/{id}", ctrl. FindByID)
		r.Put("/{id}", ctrl.Update)
		r.Delete("/{id}", ctrl.Delete)
	})
}
