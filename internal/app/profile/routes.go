package profile

import (
	"user-profiles/configs"
	"github.com/go-chi/chi/v5"
	"user-profiles/internal/middleware"
	"user-profiles/internal/infrastructure/security"
)

type ProfileHandler struct {
	*configs.Config
	Service
	security.JWTService
}

type ProfileHandlerDeps struct {
	*configs.Config
	Service
	security.JWTService
}

func NewProfileHandler(router chi.Router, deps ProfileHandlerDeps) {
	handler := &ProfileHandler {
		Config: deps.Config,
		Service: deps.Service,
		JWTService: deps.JWTService,
	}
	
	router.Route(("/profile"), func(router chi.Router) {
		router.Use(middleware.JWTMiddleware(handler.JWTService))

		router.Get("/", handler.View)
		router.Patch("/name", handler.Update)
		router.Delete("/", handler.Delete)
	})
}
