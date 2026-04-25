package app

import (
	"fmt"
	"net/http"
	"user-profiles/configs"

	"user-profiles/internal/auth"
	"user-profiles/internal/http/middleware"
	"user-profiles/internal/storage/db"
	"user-profiles/internal/users"

	"user-profiles/internal/http/handlers"
	"user-profiles/internal/http/routes"
	auth_postgres "user-profiles/internal/storage/postgres/auth"
	users_postgres "user-profiles/internal/storage/postgres/users"

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
	userRepo := users_postgres.NewUsersRepository(dbConn)
	tokenRepo := auth_postgres.NewTokenRepository(dbConn)

	// Services
	jwtService := auth.NewJWTService(conf.Secret)
	refreshTokenService := auth.NewRefreshTokenService()
	authService := auth.NewAuthService(userRepo, tokenRepo, jwtService, refreshTokenService)
	userService := users.NewUsersService(userRepo)

	// Mux
	mux := chi.NewRouter()

	// Middlewares
	mux.Use(middleware.CORS)

	auth_handler := handlers.NewAuthHandler(mux, handlers.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
		JWTService:  jwtService,
	})
	routes.InitAuthRoutes(mux, auth_handler)

	user_handler := handlers.NewUserHandler(mux, handlers.UserHandlerDeps{
		Config:       conf,
		UsersService: userService,
		AuthService:  authService,
		JWTService:   jwtService,
	})
	routes.InitUserRoutes(mux, user_handler)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("Server is listening on port 8080")
	return server.ListenAndServe()
}
