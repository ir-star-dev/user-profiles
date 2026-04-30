package routes

import (
	"user-profiles/internal/http/handlers"
	"user-profiles/internal/http/middleware"

	"github.com/go-chi/chi/v5"
)

func InitAuthRoutes(router chi.Router, handler *handlers.AuthHandler) {
	router.Route(("/auth"), func(router chi.Router) {
		router.With(middleware.CheckAuthAndRedirect(handler.JWTService)).Get("/register", handler.RegisterPage)
		router.Post("/register", handler.Register)

		router.With(middleware.CheckAuthAndRedirect(handler.JWTService)).Get("/login", handler.LoginPage)
		router.Post("/login", handler.Login)

		router.With(middleware.StrictAuthMiddleware(handler.JWTService)).Post("/logout", handler.Logout)
	})
}
