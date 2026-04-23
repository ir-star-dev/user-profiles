package routes

import (
	"user-profiles/internal/http/middleware"
	"user-profiles/internal/http/handlers"

	"github.com/go-chi/chi/v5"
)

func InitUserRoutes(router chi.Router, handler *handlers.UserHandler) {
	// Other users
	router.Route(("/profile"), func(router chi.Router) {
		router.Use(middleware.SoftAuthMiddleware(handler.JWTService, handler.AuthService))

		router.Get("/", handler.ProfilePage)

		// router.With(middleware.BanMiddleware).Patch("/name", handler.Update)
		// router.With(middleware.BanMiddleware).Delete("/", handler.Delete)
	})
}
