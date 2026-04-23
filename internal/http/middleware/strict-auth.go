package middleware

import (
	"context"
	"net/http"
	"user-profiles/internal/auth"
	"user-profiles/internal/http/cookie"
	"user-profiles/internal/http/resp"
)

type contextKey string

const (
    UserIdKey contextKey = "user_id"
    UserRoleKey contextKey = "user_role"
	UserBanKey contextKey = "user_ban"
)

func StrictAuthMiddleware(jwtService auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessCookie, err := cookie.Get("__up_access_token", r)
			if accessCookie == nil || accessCookie.Value == "" {
				resp.Json(w, "Unauthorized", http.StatusUnauthorized)
				return
			}		

			token := accessCookie.Value
			claims, err := jwtService.Parse(token)
			if err != nil {
				resp.Json(w, "Unauthorized", http.StatusUnauthorized)
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