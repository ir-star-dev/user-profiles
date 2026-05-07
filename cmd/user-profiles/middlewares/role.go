package middlewares

import (
	"net/http"
	"user-profiles/cmd/user-profiles/utils"
	"user-profiles/internal/http/resp"
)

func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	roleSet := make(map[string]struct{})
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			role, err := utils.GetUserRole(r.Context())
			if err != nil {
				resp.Json(w, err.Error(), http.StatusForbidden)
				return
			}

			if _, allowed := roleSet[role]; !allowed {
				resp.Json(w, "You are not allowed to see this page", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
