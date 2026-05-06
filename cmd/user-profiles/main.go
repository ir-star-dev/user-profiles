package main

import (
	"log"
	"net/http"
	"user-profiles/cmd/user-profiles/auth"
	"user-profiles/cmd/user-profiles/handlers"
	"user-profiles/cmd/user-profiles/posts"
	"user-profiles/cmd/user-profiles/users"
	"user-profiles/configs"
	"user-profiles/internal/http/middleware"
	"user-profiles/internal/storage/db"
	auth_postgres "user-profiles/internal/storage/postgres/auth"
	posts_postgres "user-profiles/internal/storage/postgres/posts"
	users_postgres "user-profiles/internal/storage/postgres/users"

	"github.com/go-chi/chi/v5"
)

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

	// Mux
	mux := chi.NewRouter()

	// Middlewares
	mux.Use(middleware.CORS)

	auth_handler := handlers.NewAuthHandler(mux, handlers.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
		JWTService:  jwtService,
	})
	user_handler := handlers.NewUserHandler(mux, handlers.UserHandlerDeps{
		Config:       conf,
		AuthService:  authService,
		JWTService:   jwtService,
		UsersService: userService,
	})
	post_handler := handlers.NewPostHandler(mux, handlers.PostHandlerDeps{
		Config:      conf,
		AuthService: authService,
		JWTService:  jwtService,
		PostService: postService,
	})
	InitRoutes(mux, auth_handler, user_handler, post_handler)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("Server is listening on port 8080")
	server.ListenAndServe()
}
