package main

import (
	"net/http"
	"user-profiles/cmd/user-profiles/handlers"
	"user-profiles/internal/http/middleware"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(router chi.Router, ah *handlers.AuthHandler, uh *handlers.UserHandler, ph *handlers.PostHandler) {
	fs := http.FileServer(http.Dir("././ui/static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	router.Route(("/auth"), func(router chi.Router) {
		router.With(middleware.CheckAuthAndRedirect(ah.JWTService)).Get("/register", ah.RegisterPage)
		router.Post("/register", ah.Register)

		router.With(middleware.CheckAuthAndRedirect(ah.JWTService)).Get("/login", ah.LoginPage)
		router.Post("/login", ah.Login)

		router.With(middleware.StrictAuthMiddleware(ah.JWTService)).Post("/logout", ah.Logout)
	})

	router.Get("/", ph.MainPage)
	router.Get("/posts/post/{slug}", ph.ViewPost)
	router.Get("/posts/user/{username}", ph.ViewUserPosts)

	router.Route(("/panel"), func(router chi.Router) {
		router.Use(middleware.SoftAuthMiddleware(ah.JWTService, ah.AuthService))

		// Role-limited action types on posts, role checks based on GET settings
		// ?role=admin
		// ?role=editor
		// ?role=user
		router.Get("/posts", ph.PostList)
		router.Get("/posts/post/{id}/create", ph.CreateForm)
		router.Get("/posts/post/{id}/edit", ph.EditForm)

		// After creating, the post will have a review status.
		router.Post("/posts/post/{id}/create", ph.Create)
		// After editing, the post will have a review status.
		router.Patch("/posts/post/{id}/edit", ph.Edit)

		router.With(middleware.RoleMiddleware("admin", "editor")).Patch("/posts/post/{id}/review", ph.Review)
		router.With(middleware.RoleMiddleware("admin", "editor")).Patch("/posts/post/{id}/publish", ph.Publish)
		router.With(middleware.RoleMiddleware("admin", "editor")).Delete("/posts/post/{id}", ph.Delete)

		router.With(middleware.RoleMiddleware("admin")).Get("/users/page/{page}", uh.UserListPage)

		router.Route(("/profile"), func(router chi.Router) {
			router.Use(middleware.SoftAuthMiddleware(ah.JWTService, ah.AuthService))

			router.Get("/{id}", uh.ProfilePage)

			// router.With(middleware.BanMiddleware).Patch("/{id}/name", handler.Update)
			router.With(middleware.RoleMiddleware("admin")).Patch("/{id}/ban", uh.Ban)
			router.With(middleware.RoleMiddleware("admin")).Patch("/{id}/unban", uh.Unban)

			router.With(middleware.BanMiddleware).Get("/{id}/delete-confirm", uh.DeleteConfirm)
			router.With(middleware.BanMiddleware).Delete("/{id}", uh.Delete)
		})
	})
}
