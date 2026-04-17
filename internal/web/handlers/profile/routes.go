package profile_handlers

import (
	"html/template"
	"user-profiles/configs"
	//"user-profiles/internal/middleware"
	"user-profiles/internal/security"
	profile_service "user-profiles/internal/service/profile"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Config     *configs.Config
	Service    profile_service.Service
	JWTService security.JWTService
	Tmpl       *template.Template
}

type HandlerDeps struct {
	Config     *configs.Config
	Service    profile_service.Service
	JWTService security.JWTService
	Tmpl       *template.Template
}

func New(router chi.Router, deps HandlerDeps) {
	handler := &Handler{
		Config:     deps.Config,
		Service:    deps.Service,
		JWTService: deps.JWTService,
		Tmpl:       deps.Tmpl,
	}

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

	// Other users
	router.Route(("/profile"), func(router chi.Router) {
		//router.Use(middleware.AuthMiddleware(handler.JWTService))

		router.Get("/", handler.ProfilePage)

		// router.With(middleware.BanMiddleware).Patch("/name", handler.Update)
		// router.With(middleware.BanMiddleware).Delete("/", handler.Delete)
	})

}
