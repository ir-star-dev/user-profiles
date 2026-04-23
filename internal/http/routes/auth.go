package routes

import (
	"user-profiles/internal/http/middleware"
	"user-profiles/internal/http/handlers"

	"github.com/go-chi/chi/v5"
)

func InitAuthRoutes(router chi.Router, handler *handlers.AuthHandler) {

	router.With(middleware.RedirectIfAuth(handler.JWTService)).Get("/", handler.MainPage)


	router.Route(("/auth"), func(router chi.Router) {
		router.With(middleware.RedirectIfAuth(handler.JWTService)).Get("/register", handler.RegisterPage)
		router.Post("/register", handler.Register)
		
		router.With(middleware.RedirectIfAuth(handler.JWTService)).Get("/login", handler.LoginPage)
		router.Post("/login", handler.Login)

		router.With(middleware.StrictAuthMiddleware(handler.JWTService)).Post("/logout", handler.Logout)
	})
}
