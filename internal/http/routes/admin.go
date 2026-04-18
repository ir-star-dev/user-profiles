package routes

import (
	//"user-profiles/internal/middleware"
	"user-profiles/internal/http/handlers"

	"github.com/go-chi/chi/v5"
)

func InitAdminRoutes(router chi.Router, handler *handlers.AdminHandler) {
	// Admin
	// router.Route(("/profiles"), func(router chi.Router) {
	// 	router.Use(middleware.AuthMiddleware(handler.JWTService))
	// 	router.Use(middleware.RoleMiddleware("admin"))

	// 	router.Get("/", handler.ViewAllByAdmin)
	// 	router.Get("/{id}", handler.ViewByAdmin)

	// 	router.Patch("/name/{id}", handler.UpdateByAdmin)
	// 	router.Patch("/{id}/ban", handler.BanByAdmin)
	// 	router.Patch("/{id}/unban", handler.UnbanByAdmin)

	// 	router.Delete("/{id}", handler.DeleteByAdmin)
	// })
}
