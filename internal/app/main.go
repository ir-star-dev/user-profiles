package app

import (
	"log/slog"
	"net/http"
	"os"
	"user-profiles/configs"
	"user-profiles/internal/auth"
	"user-profiles/internal/dashboard"
	"user-profiles/internal/http/middlewares"
	"user-profiles/internal/posts"
	"user-profiles/internal/storage/db"
	auth_postgres "user-profiles/internal/storage/postgres/auth"
	posts_postgres "user-profiles/internal/storage/postgres/posts"
	users_postgres "user-profiles/internal/storage/postgres/users"
	"user-profiles/internal/templates"
	"user-profiles/internal/users"

	"github.com/go-chi/chi/v5"
	"github.com/lmittmann/tint"
)

func Run() error {
	logger := slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level: slog.LevelDebug,
		}),
	)
	// Config
	conf, err := configs.Load()
	if err != nil {
		logger.Error("Failed to load config")
		return err
	}

	// DB
	dbConn, err := db.Connect(conf)
	if err != nil {
		logger.Error("Failed to connect db")
		return err
	}
	defer dbConn.Close()

	// Repositories
	userRepo := users_postgres.NewUsersRepository(dbConn)
	tokenRepo := auth_postgres.NewTokenRepository(dbConn)
	loginsRepo := auth_postgres.NewLoginsRepository(dbConn)
	postRepo := posts_postgres.NewPostRepository(dbConn)

	// Services
	jwtService := auth.NewJWTService(conf.Secret)
	rTService := auth.NewRefreshTokenService()
	aS := auth.NewAuthService(userRepo, tokenRepo, loginsRepo, jwtService, rTService)
	uS := users.NewUsersService(userRepo)
	pS := posts.NewPostService(postRepo)
	dS := dashboard.NewDashboardService(uS, pS, aS)

	tc, err := templates.NewTemplateCache()
	if err != nil {
		logger.Warn(err.Error())
		return err
	}

	// Mux
	mux := chi.NewRouter()
	// Middlewares
	mux.Use(
		middlewares.Logger(logger),
		tc.RecoverPanic,
		middlewares.CORS,
		middlewares.CSRF(false),
	)

	dH := dashboard.NewDashboardHandler(mux, dashboard.DHandlerDeps{
		UService: uS,
		PService: pS,
		TCache:   *tc,
		DService: *dS,
	})
	aH := auth.NewAuthHandler(mux, auth.AHandlerDeps{
		Config:     conf,
		AService:   aS,
		JWTService: jwtService,
		TCache:     *tc,
	})
	pH := posts.NewPostHandler(mux, posts.PHandlerDeps{
		PService: pS,
		TCache:   *tc,
	})
	uH := users.NewUserHandler(mux, users.UHandlerDeps{
		UService: uS,
		TCache:   *tc,
	})

	InitRoutes(mux, aH, pH, dH, uH, tc)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	logger.Info("Server is listening on port 8080")
	err = server.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}
