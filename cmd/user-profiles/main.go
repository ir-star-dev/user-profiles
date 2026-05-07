package main

import (
	"fmt"
	"log"
	"net/http"
	"user-profiles/cmd/user-profiles/middlewares"
	"user-profiles/configs"
	"user-profiles/internal/auth"
	"user-profiles/internal/storage/db"
	auth_postgres "user-profiles/internal/storage/postgres/auth"
	users_postgres "user-profiles/internal/storage/postgres/users"
	"user-profiles/internal/users"

	"github.com/go-chi/chi/v5"
)

type application struct {
	*configs.Config
	auth.AuthService
	auth.JWTService
	auth.RefreshTokenService
	users.UsersService
}

func main() {
	// Config
	conf, err := configs.Load()
	if err != nil {
		log.Fatalln(err)
	}

	// DB
	dbConn, err := db.Connect(conf)
	if err != nil {
		log.Fatalln(err)
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
	mux.Use(middlewares.CORS)

	app := application{
		Config:              conf,
		AuthService:         authService,
		JWTService:          jwtService,
		RefreshTokenService: refreshTokenService,
		UsersService:        userService,
	}
	app.InitRoutes(mux)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("Server is listening on port 8080")
	server.ListenAndServe()
}
