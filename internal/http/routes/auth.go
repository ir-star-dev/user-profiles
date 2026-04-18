package routes

import (
	//"user-profiles/internal/middleware"
	"user-profiles/internal/http/handlers"

	"github.com/go-chi/chi/v5"
)

func InitAuthRoutes(router chi.Router, handler *handlers.AuthHandler) {
	router.Route(("/auth"), func(router chi.Router) {

		router.Get("/register", handler.RegisterPage)
		router.Post("/register", handler.Register)

		router.Get("/login", handler.LoginPage)
		router.Post("/login", handler.Login)

		// router.With(middleware.AuthMiddleware(handler.JWTService)).Post("/logout", handler.Logout)
		// router.Post("/refresh", handler.Refresh)
	})
}
