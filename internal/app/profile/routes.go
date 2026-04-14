package profile

import (
	"user-profiles/configs"
	"github.com/go-chi/chi/v5"
	"user-profiles/internal/middleware"
	"user-profiles/internal/security"
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

	// Admin 
	router.Route(("/profiles"), func(router chi.Router) {
		router.Use(middleware.AuthMiddleware(handler.JWTService))
		router.Use(middleware.RoleMiddleware("admin"))

		router.Get("/", handler.ViewAll)

		router.Get("/{id}", handler.View)
		router.Patch("/name/{id}", handler.Update)
		router.Delete("/{id}", handler.Delete)
	})


	// Current user
	router.Route(("/profile"), func(router chi.Router) {
		router.Use(middleware.AuthMiddleware(handler.JWTService))

		router.Get("/", handler.View)
		router.Patch("/name", handler.Update)
		router.Delete("/", handler.Delete)
	})
	
}

