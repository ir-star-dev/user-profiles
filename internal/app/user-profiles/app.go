package app

import (
	"fmt"
	"html/template"
	"net/http"
	"user-profiles/configs"

	//"user-profiles/internal/app/profile"
	"user-profiles/internal/middleware"
	"user-profiles/internal/security"

	auth_service "user-profiles/internal/service/auth"
	profile_service "user-profiles/internal/service/profile"

	"user-profiles/internal/storage/db"

	postgres "user-profiles/internal/storage/postgres/user"
	auth_handlers "user-profiles/internal/web/handlers/auth"
	profile_handlers "user-profiles/internal/web/handlers/profile"

	"github.com/go-chi/chi/v5"
)

func Run() error {
	// Config
	conf, err := configs.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// DB
	dbConn, err := db.Connect(conf)
	if err != nil {
		return fmt.Errorf("failed to connect db: %w", err)
	}
	defer dbConn.Close()

	// Repositories
	userRepo := postgres.NewUserRepository(dbConn)
	tokenRepo := security.NewTokenRepository(dbConn)

	// Services
	jwtService := security.NewJWTService(conf.Secret)
	refreshTokenService := security.NewRefreshTokenService()

	authService := auth_service.New(userRepo, tokenRepo, jwtService, refreshTokenService)
	profileService := profile_service.New(userRepo)

	// Templates
	tmpl := template.Must(template.ParseGlob("./internal/web/templates/parts/*/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("./internal/web/templates/pages/*.html"))

	// Handlers
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.CORS)

	auth_handlers.New(r, auth_handlers.HandlerDeps{
		Config:     conf,
		Service:    authService,
		JWTService: jwtService,
		Tmpl:       tmpl,
	})

	profile_handlers.New(r, profile_handlers.HandlerDeps{
		Config:     conf,
		Service:    profileService,
		JWTService: jwtService,
		Tmpl:       tmpl,
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	fmt.Println("Server is listening on port 8080")

	return server.ListenAndServe()
}
