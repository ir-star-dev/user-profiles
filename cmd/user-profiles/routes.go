package main

import (
	"net/http"
	"user-profiles/cmd/user-profiles/middlewares"

	"github.com/go-chi/chi/v5"
)

func (app *application) InitRoutes(mux chi.Router) {

	fs := http.FileServer(http.Dir("././ui/static"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fs))

	mux.With(middlewares.SoftAuthMiddleware(app.JWTService, app.AuthService)).Get("/", app.Index)

	mux.Route(("/auth"), func(m chi.Router) {
		m.With(middlewares.CheckAuthAndRedirect(app.JWTService)).Get("/register", app.RegisterForm)
		m.With(middlewares.CheckAuthAndRedirect(app.JWTService)).Get("/login", app.LoginForm)
		
		m.Post("/register", app.Register)
		m.Post("/login", app.Login)
		m.With(middlewares.StrictAuthMiddleware(app.JWTService)).Post("/logout", app.Logout)
	})

	mux.Route(("/profile"), func(m chi.Router) {
		m.Use(middlewares.SoftAuthMiddleware(app.JWTService, app.AuthService))

		m.Get("/{id}", app.Profile)
		m.With(middlewares.RoleMiddleware("admin")).Get("/users/page/{page}", app.UserList)

		// router.With(middleware.BanMiddleware).Patch("/{id}/name", handler.Update)
		m.With(middlewares.RoleMiddleware("admin")).Patch("/{id}/ban", app.Ban)
		m.With(middlewares.RoleMiddleware("admin")).Patch("/{id}/unban", app.Unban)

		m.With(middlewares.BanMiddleware).Get("/{id}/delete-confirm", app.DeleteConfirm)
		m.With(middlewares.BanMiddleware).Delete("/{id}", app.Delete)
	})
}
