package app

import (
	"fmt"
	"net/http"
	"user-profiles/configs"
	"user-profiles/internal/app/auth"
	"user-profiles/internal/app/profile"
	"user-profiles/internal/middleware"
	"user-profiles/internal/infrastructure/db/postgres/user"
	"user-profiles/internal/infrastructure/security"
	"user-profiles/pkg/db"
	//"user-profiles/pkg/logger"

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

	authService := auth.NewAuthService(userRepo, tokenRepo, jwtService, refreshTokenService)
	profileService := profile.NewProfileService(userRepo)

	// Logger 
	//log := logger.New()

	// Handlers
	r := chi.NewRouter()
	// Middlewares
	//r.Use(middleware.Logger(log))
	r.Use(middleware.CORS)

	auth.NewAuthHandler(r, auth.AuthHandlerDeps{
		Config:      conf,
		Service: authService,
		JWTService: jwtService,
	})
	profile.NewProfileHandler(r, profile.ProfileHandlerDeps{
		Config:         conf,
		Service: profileService,
		JWTService: jwtService,
	})
	
	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	fmt.Println("Server is listening on port 8080")
	server.ListenAndServe()

	return nil
}
