package main

import (
	"net/http"
	"user-profiles/cmd/user-profiles/handlers"
	"user-profiles/cmd/user-profiles/middlewares"
	"user-profiles/cmd/user-profiles/panel"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(router chi.Router, ah *handlers.AuthHandler, ph *handlers.PostHandler, dh *handlers.DashboardHandler, tc *panel.Templates) {
	fs := http.FileServer(http.Dir("././ui/static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	router.Route(("/auth"), func(router chi.Router) {
		router.With(middlewares.CheckAuthAndRedirect(ah.JWTService)).Get("/signup", ah.SignupForm)
		router.Post("/signup", ah.Signup)

		router.With(middlewares.CheckAuthAndRedirect(ah.JWTService)).Get("/login", ah.LoginForm)
		router.Post("/login", ah.Login)

		router.With(middlewares.StrictAuthMiddleware(ah.JWTService, ah.AuthService)).Post("/logout", ah.Logout)
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		tc.NotFound(w, r)
	})

	router.With(middlewares.SoftAuthMiddleware(ah.JWTService, ah.AuthService)).Get("/", ph.Home)
	router.With(middlewares.SoftAuthMiddleware(ah.JWTService, ah.AuthService)).Get("/posts/post/{slug}", ph.ViewPost)
	router.With(middlewares.SoftAuthMiddleware(ah.JWTService, ah.AuthService)).Get("/posts/user/{username}", ph.ViewUserPosts)

	router.Route(("/panel"), func(router chi.Router) {
		router.NotFound(func(w http.ResponseWriter, r *http.Request) {
			tc.PanelNotFound(w, r)
		})
		router.Use(middlewares.SoftAuthMiddleware(ah.JWTService, ah.AuthService))

		router.Get("/posts", dh.Posts)
		// router.Get("/post/create", ph.CreateForm)
		// router.Get("/posts/post/{id}/edit", ph.EditForm)

		// After creating, the post will have a review status.
		//router.Post("/post/create", ph.Create)
		// After editing, the post will have a review status.
		//router.Patch("/posts/post/{id}/edit", ph.Edit)

		// router.With(middlewares.RoleMiddleware("admin", "editor")).Patch("/posts/post/{id}/review", ph.Review)
		// router.With(middlewares.RoleMiddleware("admin", "editor")).Patch("/posts/post/{id}/publish", ph.Publish)
		// router.With(middlewares.RoleMiddleware("admin", "editor")).Get("/posts/post/{id}/delete-confirm", ph.DeleteConfirm)
		// router.With(middlewares.RoleMiddleware("admin", "editor")).Delete("/posts/post/{id}", ph.Delete)

		router.With(middlewares.RoleMiddleware("admin")).Get("/users", dh.Users)

		router.Route(("/profile"), func(router chi.Router) {
			router.Use(middlewares.SoftAuthMiddleware(ah.JWTService, ah.AuthService))

			router.Get("/{id}", dh.Profile)

			router.With(middlewares.BanMiddleware).Get("/{id}/name", dh.UpdateNameModal)
			router.With(middlewares.BanMiddleware).Patch("/{id}/name", dh.UpdateName)
			router.With(middlewares.RoleMiddleware("admin")).Patch("/{id}/ban", dh.Ban)
			router.With(middlewares.RoleMiddleware("admin")).Patch("/{id}/unban", dh.Unban)
			router.With(middlewares.BanMiddleware).Get("/{id}/delete-confirm", dh.DeleteConfirm)
			router.With(middlewares.BanMiddleware).Delete("/{id}", dh.Delete)
		})
	})
}
