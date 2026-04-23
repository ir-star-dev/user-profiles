package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
	"user-profiles/internal/auth"
	"user-profiles/internal/http/cookie"
)

func SoftAuthMiddleware(jwtService auth.JWTService, authService auth.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessCookie, _ := cookie.Get("__up_access_token", r)
			refreshCookie, err := cookie.Get("__up_refresh_token", r)
			if accessCookie == nil {
				if refreshCookie != nil && refreshCookie.Value != ""  {
					tokens, err := authService.Refresh(refreshCookie.Value)
					log.Println(err.Error())
					if err == nil {
						cookie.Set(tokens.Access, "__up_access_token", 5*time.Minute, w)
						cookie.Set(tokens.Refresh, "__up_refresh_token", 7*24*time.Hour, w)
					} else {
						http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
						return
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
					} else {
						http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
						return
					}
				} else {
					http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
					return
				}
			}

			sub, ok := claims["sub"].(float64)
			if !ok {
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}
			uId := int(sub)
			// send UserIdKey in context
			ctx := context.WithValue(r.Context(), UserIdKey, uId)

			role, ok := claims["role"].(string)
			if !ok {
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}
			// send UserRoleKey in context
			ctx = context.WithValue(ctx, UserRoleKey, role)
			ban, ok := claims["banned"].(bool)
			if !ok {
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}
			// send UseBanKey in context
			ctx = context.WithValue(ctx, UserBanKey, ban)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
