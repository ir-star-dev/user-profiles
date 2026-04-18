package middleware

import (
	"context"
	"net/http"
	"strings"
	"user-profiles/internal/http/resp"
	"user-profiles/internal/auth"
)

type contextKey string

const (
    UserIdKey contextKey = "user_id"
    UserRoleKey contextKey = "user_role"
	UserBanKey contextKey = "user_ban"
)

func AuthMiddleware(jwtService auth.JWTService) func(http.Handler) http.Handler {
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

			sub, ok := claims["sub"].(float64)
			if !ok {
				resp.Json(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			uId := int(sub)
			// send UserIdKey in context
			ctx := context.WithValue(r.Context(), UserIdKey, uId)

			role, ok := claims["role"].(string)
			if !ok {
				resp.Json(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			// send UserRoleKey in context		
			ctx = context.WithValue(ctx, UserRoleKey, role)
			ban, ok := claims["banned"].(bool)
			if !ok {
				resp.Json(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			// send UseBanKey in context
			ctx = context.WithValue(ctx, UserBanKey, ban)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}