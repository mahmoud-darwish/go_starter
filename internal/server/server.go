package server

import (
	"net/http"

	"starter/config"
	"starter/internal/middleware"
	productRoutes "starter/product/routes"
	userRoutes "starter/user/routes"
	notificationsRoutes "starter/Notifications/routes" 
	commentsRoutes "starter/Comments/routes"
	subscriptionsRoutes "starter/Subscription/routes"
	likesRoutes "starter/Likes/routes"


	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func NewServer(db *gorm.DB) *http.Server {
	r := chi.NewRouter()

	r.Use(middleware.SetupMiddleware)

	userRoutes.RegisterUserRoutes(r, db)
	productRoutes.RegisterProductRoutes(r, db)
	notificationsRoutes.RegisterNotificationRoutes(r, db) 
	commentsRoutes.RegisterCommentsRoutes(r, db) 
	subscriptionsRoutes.RegisterSubscriptionRoutes(r, db)
	likesRoutes.RegisterLikesRoutes(r, db)

	return &http.Server{
		Addr:    ":" + config.GetConfig().Port,
		Handler: r,
	}
}
