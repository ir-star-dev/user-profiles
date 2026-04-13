package auth

import (
	"user-profiles/configs"
	"user-profiles/internal/infrastructure/security"
	"user-profiles/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type AuthHandlerDeps struct {
	*configs.Config
	Service
	security.JWTService
}

type AuthHandler struct {
	*configs.Config
	Service
	security.JWTService
}

func NewAuthHandler(router chi.Router, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:  deps.Config,
		Service: deps.Service,
		JWTService: deps.JWTService,
	}
	router.Route(("/auth"), func(router chi.Router) {
		router.Post("/login", handler.Login)
		router.Post("/register", handler.Register)

		router.With(middleware.JWTMiddleware(handler.JWTService)).Post("/logout", handler.Logout)
		router.Post("/refresh", handler.Refresh)
	})
	
}
