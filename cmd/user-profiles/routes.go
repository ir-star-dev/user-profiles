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
		router.With(middlewares.RoleMiddleware("admin")).Get("/generate-password", dh.GeneratePassword)

		router.Get("/posts", dh.Posts)

		// router.Get("/post/create", ph.CreatePostForm)
		// router.Post("/post/create", ph.CreatePost)

		// router.Get("/posts/post/{id}/edit", ph.EditPostForm)
		// router.Patch("/posts/post/{id}/edit", ph.EditPost)

		router.With(middlewares.RoleMiddleware("admin", "moderator")).Patch("/posts/post/{id}/review", dh.Review)
		router.With(middlewares.RoleMiddleware("admin", "moderator")).Patch("/posts/post/{id}/publish", dh.Publish)

		router.With(middlewares.RoleMiddleware("admin", "moderator")).Get("/posts/post/{id}/delete-confirm", dh.DeletePostConfirm)
		router.With(middlewares.RoleMiddleware("admin", "moderator")).Delete("/posts/post/{id}", dh.DeletePost)

		router.With(middlewares.RoleMiddleware("admin")).Get("/users", dh.Users)

		
		router.With(middlewares.RoleMiddleware("admin")).Get("/user/add", dh.CreateUserForm)
		router.With(middlewares.RoleMiddleware("admin")).Post("/user/add", dh.CreateUser)


		router.Route(("/profile"), func(router chi.Router) {
			router.Use(middlewares.SoftAuthMiddleware(ah.JWTService, ah.AuthService))

			router.Get("/{id}", dh.Profile)

			router.With(middlewares.BanMiddleware).Get("/{id}/name", dh.UpdateNameModal)
			router.With(middlewares.BanMiddleware).Patch("/{id}/name", dh.UpdateName)

			router.With(middlewares.RoleMiddleware("admin")).Patch("/{id}/ban", dh.Ban)
			router.With(middlewares.RoleMiddleware("admin")).Patch("/{id}/unban", dh.Unban)

			router.With(middlewares.BanMiddleware).Get("/{id}/delete-confirm", dh.DeleteUserConfirm)
			router.With(middlewares.BanMiddleware).Delete("/{id}", dh.DeleteUser)
		})
	})
}
