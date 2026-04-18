package middleware

import (
	"net/http"
	"user-profiles/internal/http/resp"
)

func BanMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ban, err := GetUserBanStatus(r.Context())
		if err != nil {
			resp.Json(w, err.Error(), http.StatusForbidden)
			return
		}
		if *ban {
			resp.Json(w, "You banned", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
