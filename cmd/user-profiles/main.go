package main

import (
	"log"
	"net/http"
	"user-profiles/cmd/user-profiles/auth"
	"user-profiles/cmd/user-profiles/handlers"
	"user-profiles/cmd/user-profiles/middlewares"
	"user-profiles/cmd/user-profiles/panel"
	"user-profiles/cmd/user-profiles/posts"
	"user-profiles/cmd/user-profiles/users"
	"user-profiles/configs"
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
	dashboard := panel.NewDashboardService(userRepo, postRepo)
	jwtService := auth.NewJWTService(conf.Secret)
	refreshTokenService := auth.NewRefreshTokenService()
	authService := auth.NewAuthService(userRepo, tokenRepo, jwtService, refreshTokenService)
	userService := users.NewUsersService(userRepo)
	postService := posts.NewPostService(postRepo)

	// Mux
	mux := chi.NewRouter()

	// Middlewares
	mux.Use(middlewares.CORS, middlewares.RecoverPanic)

	tc, err := panel.NewTemplateCache()
	if err != nil {
		log.Fatal(err.Error())
	}

	authH := handlers.NewAuthHandler(mux, handlers.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
		JWTService:  jwtService,
		TCache:      *tc,
	})

	dashboardH := handlers.NewDashboardHandler(mux, handlers.DashboardHandlerDeps{
		UsersService: userService,
		PostsService: postService,
		TCache:       *tc,
		Dashboard:    *dashboard,
	})

	postH := handlers.NewPostHandler(mux, handlers.PostHandlerDeps{
		PostService: postService,
		TCache:      *tc,
	})
	InitRoutes(mux, authH, postH, dashboardH, tc)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server is listening on port 8080")
	server.ListenAndServe()
}
