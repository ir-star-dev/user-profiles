package middlewares

import (
	"net/http"
	"user-profiles/internal/http/request"
	"user-profiles/internal/templates"
)

func RoleMiddleware(tc *templates.Templates, allowedRoles ...string) func(http.Handler) http.Handler {
	roleSet := make(map[string]struct{})
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, err := request.GetUserRole(r.Context())
			if err != nil {
				// HTMX request
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", "/403")
					w.WriteHeader(http.StatusForbidden)
					return
				}
				// Normal request
				tc.Forbidden(w, r)
				return
			}
			if _, allowed := roleSet[role]; !allowed {
				// HTMX request
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", "/403")
					w.WriteHeader(http.StatusForbidden)
					return
				}
				// Normal request
				tc.Forbidden(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
