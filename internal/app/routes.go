package app

import (
	"net/http"
	"user-profiles/internal/auth"
	"user-profiles/internal/dashboard"
	"user-profiles/internal/http/middlewares"
	"user-profiles/internal/posts"
	"user-profiles/internal/templates"
	"user-profiles/internal/users"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(router chi.Router, ah *auth.AHandler, ph *posts.PHandler, dh *dashboard.DHandler, uh *users.UHandler, tc *templates.Templates) {
	fs := http.FileServer(http.Dir("././ui/static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Errors
	router.Get("/403", func(w http.ResponseWriter, r *http.Request) {
		tc.Forbidden(w, r)
	})
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		tc.NotFound(w, r)
	})

	// Authentication
	router.Route(("/auth"), func(router chi.Router) {
		router.With(middlewares.CheckAuthAndRedirect(ah.JWTService)).Get("/signup", ah.SignupForm)
		router.Post("/signup", ah.Signup)

		router.With(middlewares.CheckAuthAndRedirect(ah.JWTService)).Get("/login", ah.LoginForm)
		router.Post("/login", ah.Login)

		router.With(middlewares.StrictAuthMiddleware(ah.JWTService, ah.AService)).Post("/logout", ah.Logout)
	})

	// Front
	router.With(middlewares.SoftAuthMiddleware(ah.JWTService, ah.AService)).Get("/", ph.Home)
	router.With(middlewares.SoftAuthMiddleware(ah.JWTService, ah.AService)).Get("/posts/post/{slug}", ph.ViewPost)
	router.With(middlewares.SoftAuthMiddleware(ah.JWTService, ah.AService)).Get("/posts/user/{username}", ph.ViewUserPosts)

	//Admin
	router.Route(("/panel"), func(router chi.Router) {
		router.Use(middlewares.SoftAuthMiddleware(ah.JWTService, ah.AService))

		// Dashboard
		router.With(middlewares.RoleMiddleware(tc, "admin")).Get("/dashboard", dh.Statistics)

		// Posts
		router.Get("/posts", ph.Posts)
		router.Get("/posts/post/{id}/preview", ph.PreviewPost)

		// Create post
		router.With(middlewares.BanMiddleware(tc)).Get("/post/create", ph.CreatePostForm)
		router.With(middlewares.BanMiddleware(tc)).Post("/post/create", ph.CreatePost)

		// Edit post
		router.With(middlewares.BanMiddleware(tc)).Get("/posts/post/{id}/edit", ph.EditPostForm)
		router.With(middlewares.BanMiddleware(tc)).Patch("/posts/post/{id}/edit", ph.EditPost)

		// Publish & Review post
		router.With(middlewares.BanMiddleware(tc)).With(middlewares.RoleMiddleware(tc, "admin", "moderator")).Patch("/posts/post/{id}/review", ph.Review)
		router.With(middlewares.BanMiddleware(tc)).With(middlewares.RoleMiddleware(tc, "admin", "moderator")).Patch("/posts/post/{id}/publish", ph.Publish)
		// Delete post
		router.With(middlewares.BanMiddleware(tc)).With(middlewares.RoleMiddleware(tc, "admin", "moderator")).Get("/posts/post/{id}/delete-confirm", ph.DeletePostConfirm)
		router.With(middlewares.BanMiddleware(tc)).With(middlewares.RoleMiddleware(tc, "admin", "moderator")).Delete("/posts/post/{id}", ph.DeletePost)

		// Users
		router.With(middlewares.RoleMiddleware(tc, "admin")).Get("/users", uh.Users)

		// Create user
		router.With(middlewares.RoleMiddleware(tc, "admin")).Get("/generate-password", uh.GeneratePassword)
		router.With(middlewares.BanMiddleware(tc)).With(middlewares.RoleMiddleware(tc, "admin")).Get("/user/add", uh.CreateUserForm)
		router.With(middlewares.BanMiddleware(tc)).With(middlewares.RoleMiddleware(tc, "admin")).Post("/user/add", uh.CreateUser)

		// Profile
		router.Route(("/profile"), func(router chi.Router) {
			router.Use(middlewares.SoftAuthMiddleware(ah.JWTService, ah.AService))

			router.Get("/{id}", dh.Profile)

			// Update name
			router.With(middlewares.BanMiddleware(tc)).Get("/{id}/name", uh.UpdateNameModal)
			router.With(middlewares.BanMiddleware(tc)).Patch("/{id}/name", uh.UpdateName)

			// Ban & Unban
			router.With(middlewares.RoleMiddleware(tc, "admin")).Patch("/{id}/ban", uh.Ban)
			router.With(middlewares.RoleMiddleware(tc, "admin")).Patch("/{id}/unban", uh.Unban)

			// Delete profile
			router.With(middlewares.BanMiddleware(tc)).Get("/{id}/delete-confirm", uh.DeleteUserConfirm)
			router.With(middlewares.BanMiddleware(tc)).Delete("/{id}", uh.DeleteUser)
		})
	})
}
