package middleware

import (
	"net/http"
	"user-profiles/internal/auth"
	"user-profiles/internal/http/cookie"
)

func RedirectIfAuth(jwtService auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := cookie.Get("__up_access_token", r)
			if err == nil && cookie.Value != "" {

				_, err := jwtService.Parse(cookie.Value)
				if err == nil {
					http.Redirect(w, r, "/profile", http.StatusSeeOther)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}