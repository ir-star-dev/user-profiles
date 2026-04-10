package middleware

import (
	"context"
	"net/http"
	"strings"
	"user-profiles/internal/infrastructure/security"
	"user-profiles/pkg/resp"
)

type contextKey string

var UserIdKey contextKey = "user_id"

func JWTMiddleware(jwtService security.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				resp.Json(w, "Missing token", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				resp.Json(w, "Invalid token format", http.StatusUnauthorized)
				return
			}

			token := parts[1]
			claims, err := jwtService.Parse(token)
			if err != nil {
				resp.Json(w, err.Error(), http.StatusUnauthorized)
				return
			}

			uId := claims["sub"]
			// send UserIdKey in context
			ctx := context.WithValue(r.Context(), UserIdKey, uId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}