package auth_handlers

import (
	"html/template"
	"user-profiles/configs"

	//"user-profiles/internal/middleware"
	"user-profiles/internal/security"
	auth_service "user-profiles/internal/service/auth"

	"github.com/go-chi/chi/v5"
)

type HandlerDeps struct {
	Config     *configs.Config
	Service    auth_service.Service
	JWTService security.JWTService
	Tmpl       *template.Template
}

type Handler struct {
	Config     *configs.Config
	Service    auth_service.Service
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
	router.Route(("/auth"), func(router chi.Router) {
		
		router.Get("/register", handler.RegisterPage)
		router.Post("/register", handler.Register)

		router.Get("/login", handler.LoginPage)
		router.Post("/login", handler.Login)

		// router.With(middleware.AuthMiddleware(handler.JWTService)).Post("/logout", handler.Logout)
		// router.Post("/refresh", handler.Refresh)
	})

}
