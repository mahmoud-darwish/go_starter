package routes

import (
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"starter/Subscription/controllers"
	
	

)

func RegisterSubscriptionRoutes(r *chi.Mux, db *gorm.DB) {
	ctrl := controllers.NewSubscriptionController(db)

	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", ctrl.Create)
		r.Get("/", ctrl.FindAll)  
		r.Get("/{id}", ctrl. FindByID)
		r.Delete("/{id}", ctrl.Delete)
	})
}
