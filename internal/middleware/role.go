package middleware

import (
	"net/http"
	"user-profiles/internal/http/resp"
)

func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	roleSet := make(map[string]struct{})
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			role, err := GetUserRole(r.Context())
			if err != nil {
				resp.Json(w, "Forbidden", http.StatusForbidden)
				return
			}

			if _, allowed := roleSet[role]; !allowed {
				resp.Json(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}