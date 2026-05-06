package middleware

import (
	"context"
	"net/http"
	"time"
	"user-profiles/cmd/user-profiles/auth"
	"user-profiles/internal/http/cookie"
)

func SoftAuthMiddleware(jwtService auth.JWTService, authService auth.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var accessToken string
			accessCookie, _ := cookie.Get("__up_access_token", r)
			refreshCookie, err := cookie.Get("__up_refresh_token", r)
			if accessCookie == "" {
				if refreshCookie != ""  {
					tokens, _ := authService.Refresh(refreshCookie)
					if tokens == nil {
						cookie.Set("", "__up_access_token", -time.Minute, w)
						cookie.Set("", "__up_refresh_token", -time.Minute, w)
						http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
						return
					} else {
						accessToken = tokens.Access
						cookie.Set(tokens.Access, "__up_access_token", 5*time.Minute, w)
						cookie.Set(tokens.Refresh, "__up_refresh_token", 7*24*time.Hour, w)
					}
				}
			} 
			
			accessToken = accessCookie
			claims, err := jwtService.Parse(accessToken)
			if err != nil {
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
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
