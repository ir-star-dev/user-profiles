package main

import (
	"net/http"
	"user-profiles/cmd/user-profiles/middlewares"
	"user-profiles/cmd/user-profiles/handlers"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(router chi.Router, ah *handlers.AuthHandler, ph *handlers.PostHandler,
	//, uh *handlers.UserHandler, ph *handlers.PostHandler
	) {
	fs := http.FileServer(http.Dir("././ui/static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	router.Route(("/auth"), func(router chi.Router) {
		router.With(middlewares.CheckAuthAndRedirect(ah.JWTService)).Get("/signup", ah.SignupForm)
		router.Post("/signup", ah.Signup)

		router.With(middlewares.CheckAuthAndRedirect(ah.JWTService)).Get("/login", ah.LoginForm)
		router.Post("/login", ah.Login)

		router.With(middlewares.StrictAuthMiddleware(ah.JWTService, ah.AuthService)).Post("/logout", ah.Logout)
	})

	router.With(middlewares.SoftAuthMiddleware(ah.JWTService, ah.AuthService)).Get("/", ph.Home)
	router.With(middlewares.SoftAuthMiddleware(ah.JWTService, ah.AuthService)).Get("/posts/post/{slug}", ph.ViewPost)
	router.With(middlewares.SoftAuthMiddleware(ah.JWTService, ah.AuthService)).Get("/posts/user/{username}", ph.ViewUserPosts)

	router.Route(("/panel"), func(router chi.Router) {
		//router.Use(middlewares.SoftAuthMiddleware(ah.JWTService, ah.AuthService))

		// Role-limited action types on posts, role checks based on GET settings
		// ?role=admin
		// ?role=editor
		// ?role=user
		// router.Get("/posts", ph.PostList)
		// router.Get("/posts/post/{id}/create", ph.CreateForm)
		// router.Get("/posts/post/{id}/edit", ph.EditForm)

		// After creating, the post will have a review status.
		//router.Post("/posts/post/{id}/create", ph.Create)
		// After editing, the post will have a review status.
		//router.Patch("/posts/post/{id}/edit", ph.Edit)

		// router.With(middlewares.RoleMiddleware("admin", "editor")).Patch("/posts/post/{id}/review", ph.Review)
		// router.With(middlewares.RoleMiddleware("admin", "editor")).Patch("/posts/post/{id}/publish", ph.Publish)
		// router.With(middlewares.RoleMiddleware("admin", "editor")).Delete("/posts/post/{id}", ph.Delete)

		// router.With(middlewares.RoleMiddleware("admin")).Get("/users/page/{page}", uh.UserListPage)

		//router.Route(("/profile"), func(router chi.Router) {
			// router.Use(middlewares.SoftAuthMiddleware(ah.JWTService, ah.AuthService))

			// router.Get("/{id}", uh.ProfilePage)

			// router.With(middlewares.BanMiddleware).Patch("/{id}/name", handler.Update)
		// 	router.With(middlewares.RoleMiddleware("admin")).Patch("/{id}/ban", uh.Ban)
		// 	router.With(middlewares.RoleMiddleware("admin")).Patch("/{id}/unban", uh.Unban)

		// 	router.With(middlewares.BanMiddleware).Get("/{id}/delete-confirm", uh.DeleteConfirm)
		// 	router.With(middlewares.BanMiddleware).Delete("/{id}", uh.Delete)
		// })
	})
}
