package app

import (
	"fmt"
	"html/template"
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
	adminService := users.NewAdminService(userRepo)

	// Templates
	tmpl := template.Must(template.ParseGlob("./internal/web/templates/parts/*/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("./internal/web/templates/pages/*.html"))

	// Mux
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.CORS)

	auth_handler := handlers.NewAuthHandler(r, handlers.AuthHandlerDeps{
		Config:     conf,
		Service:    authService,
		JWTService: jwtService,
		Tmpl:       tmpl,
	})
	routes.InitAuthRoutes(r, auth_handler)

	admin_handler := handlers.NewAdminHandler(r, handlers.AdminHandlerDeps{
		Config:     conf,
		Service:    adminService,
		JWTService: jwtService,
		Tmpl:       tmpl,
	})
	routes.InitAdminRoutes(r, admin_handler)

	user_handler := handlers.NewUserHandler(r, handlers.UserHandlerDeps{
		Config:     conf,
		Service:    userService,
		JWTService: jwtService,
		Tmpl:       tmpl,
	})
	routes.InitUserRoutes(r, user_handler)

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	fmt.Println("Server is listening on port 8080")
	return server.ListenAndServe()
}
