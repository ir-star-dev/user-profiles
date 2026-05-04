package routes

import (
	"user-profiles/internal/http/middleware"
	"user-profiles/internal/http/handlers"

	"github.com/go-chi/chi/v5"
)

func InitUserRoutes(router chi.Router, handler *handlers.UserHandler) {
	router.Get("/", handler.MainPage)

	router.Route(("/profile"), func(router chi.Router) {
		router.Use(middleware.SoftAuthMiddleware(handler.JWTService, handler.AuthService))

		router.Get("/{id}", handler.ProfilePage)

		// router.With(middleware.BanMiddleware).Patch("/{id}/name", handler.Update)
		router.With(middleware.RoleMiddleware("admin")).Patch("/{id}/ban", handler.Ban)
		router.With(middleware.RoleMiddleware("admin")).Patch("/{id}/unban", handler.Unban)

		router.With(middleware.BanMiddleware).Get("/{id}/delete-confirm", handler.DeleteConfirm)
		router.With(middleware.BanMiddleware).Delete("/{id}", handler.Delete)

		router.With(middleware.RoleMiddleware("admin")).Get("/users/page/{page}", handler.UserListPage)
	})
}
