package middleware

import (
	"context"
	"net/http"
	"time"
	"user-profiles/internal/auth"
	"user-profiles/internal/http/cookie"
	"user-profiles/internal/http/resp"
)

func SoftAuthMiddleware(jwtService auth.JWTService, authService auth.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessCookie, err := cookie.Get("__up_access_token", r)
			refreshCookie, _ := cookie.Get("__up_refresh_token", r)

			if accessCookie == nil || accessCookie.Value == "" {
				if refreshCookie != nil && refreshCookie.Value != "" {
					tokens, err := authService.Refresh(refreshCookie.Value)
					if err == nil {
						cookie.Set(tokens.Access, "__up_access_token", 5*time.Minute, w)
						cookie.Set(tokens.Refresh, "__up_refresh_token", 7*24*time.Hour, w)
					}
				}
			}

			token := accessCookie.Value
			claims, err := jwtService.Parse(token)
			if err != nil {
				if refreshCookie != nil {

					tokens, err := authService.Refresh(refreshCookie.Value)
					if err == nil {
						cookie.Set(tokens.Access, "__up_access_token", 5*time.Minute, w)
						cookie.Set(tokens.Refresh, "__up_refresh_token", 7*24*time.Hour, w)
						claims, err = jwtService.Parse(tokens.Access)
						if err != nil {
							resp.Json(w, "Unauthorized", http.StatusUnauthorized)
							return
						}
					} else {
						resp.Json(w, "Unauthorized", http.StatusUnauthorized)
						return
					}
				} else {
					resp.Json(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
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
