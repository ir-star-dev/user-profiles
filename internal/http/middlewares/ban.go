package middlewares

import (
	"net/http"
	"user-profiles/internal/http/request"
	"user-profiles/internal/templates"
)

func BanMiddleware(tc *templates.Templates) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ban, err := request.GetUserBanStatus(r.Context())
			if err != nil || *ban {
				// HTMX request
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", "/403")
					w.WriteHeader(http.StatusForbidden)
					return
				}
				// Normal browser request
				tc.Forbidden(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}