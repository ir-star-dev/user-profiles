package main

import (
	"log"
	"net/http"
	"user-profiles/cmd/user-profiles/handlers"
	"user-profiles/configs"
	"user-profiles/internal/auth"
	"user-profiles/internal/http/middleware"
	"user-profiles/internal/posts"
	"user-profiles/internal/storage/db"
	auth_postgres "user-profiles/internal/storage/postgres/auth"
	posts_postgres "user-profiles/internal/storage/postgres/posts"
	users_postgres "user-profiles/internal/storage/postgres/users"
	"user-profiles/internal/users"

	"github.com/go-chi/chi/v5"
)

type services struct {
	auth.AuthService
	auth.JWTService
	auth.RefreshTokenService
	users.UsersService
	posts.PostService
}

type app struct {
	*configs.Config
	*services
}

func main() {
	// Config
	conf, err := configs.Load()
	if err != nil {
		log.Fatal("Failed to load config: %w", err)
	}

	// DB
	dbConn, err := db.Connect(conf)
	if err != nil {
		log.Fatal("Failed to connect db: %w", err)
	}
	defer dbConn.Close()

	// Repositories
	userRepo := users_postgres.NewUsersRepository(dbConn)
	tokenRepo := auth_postgres.NewTokenRepository(dbConn)
	postRepo := posts_postgres.NewPostRepository(dbConn)

	// Services
	jwtService := auth.NewJWTService(conf.Secret)
	refreshTokenService := auth.NewRefreshTokenService()
	authService := auth.NewAuthService(userRepo, tokenRepo, jwtService, refreshTokenService)
	userService := users.NewUsersService(userRepo)
	postService := posts.PostService(postRepo)

	services := services{
		authService,
		jwtService,
		refreshTokenService,
		userService,
		postService,
	}
	app := app{
		conf,
		&services,
	}
	// Mux
	mux := chi.NewRouter()

	// Middlewares
	mux.Use(middleware.CORS)

	auth_handler := handlers.NewAuthHandler(mux, handlers.AuthHandlerDeps{
		Config:      app.Config,
		AuthService: app.services.AuthService,
		JWTService:  app.services.JWTService,
	})
	user_handler := handlers.NewUserHandler(mux, handlers.UserHandlerDeps{
		Config:       app.Config,
		AuthService:  app.services.AuthService,
		JWTService:   app.services.JWTService,
		UsersService: app.services.UsersService,
	})
	post_handler := handlers.NewPostHandler(mux, handlers.PostHandlerDeps{
		Config:      app.Config,
		AuthService: app.services.AuthService,
		JWTService:  app.services.JWTService,
		PostService: app.services.PostService,
	})
	InitRoutes(mux, auth_handler, user_handler, post_handler)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("Server is listening on port 8080")
	server.ListenAndServe()
}
