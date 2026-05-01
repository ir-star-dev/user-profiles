package routes

import (
	"net/http"
	"user-profiles/internal/http/handlers"
	"user-profiles/internal/http/middleware"

	"github.com/go-chi/chi/v5"
)

func InitAuthRoutes(router chi.Router, handler *handlers.AuthHandler) {

	fs := http.FileServer(http.Dir("././ui/static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	router.Route(("/auth"), func(router chi.Router) {
		router.With(middleware.CheckAuthAndRedirect(handler.JWTService)).Get("/register", handler.RegisterPage)
		router.Post("/register", handler.Register)

		router.With(middleware.CheckAuthAndRedirect(handler.JWTService)).Get("/login", handler.LoginPage)
		router.Post("/login", handler.Login)

		router.With(middleware.StrictAuthMiddleware(handler.JWTService)).Post("/logout", handler.Logout)
	})
}
