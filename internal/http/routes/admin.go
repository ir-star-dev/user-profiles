package routes

import (
	"user-profiles/internal/http/middleware"
	"user-profiles/internal/http/handlers"

	"github.com/go-chi/chi/v5"
)

func InitAdminRoutes(router chi.Router, handler *handlers.AdminHandler) {
	// Admin
	router.Route(("/profiles"), func(router chi.Router) {
		router.Use(middleware.AuthMiddleware(handler.JWTService))
		router.Use(middleware.RoleMiddleware("admin"))

		router.Get("/profile", handler.AdminProfilePage)

		// router.Patch("/profile/name/{id}", handler.UpdateByAdmin)
		// router.Patch("/profile/{id}/ban", handler.BanByAdmin)
		// router.Patch("/profile/{id}/unban", handler.UnbanByAdmin)

		// router.Delete("/profile/{id}", handler.DeleteByAdmin)
	})
}
