package middlewares

import (
	"net/http"
	"user-profiles/cmd/user-profiles/panel"
	"user-profiles/cmd/user-profiles/utils"
)

func BanMiddleware(tc *panel.Templates) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ban, err := utils.GetUserBanStatus(r.Context())
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