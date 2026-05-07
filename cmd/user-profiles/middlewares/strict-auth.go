package middlewares

import (
	"context"
	"net/http"
	"user-profiles/cmd/user-profiles/utils"
	"user-profiles/internal/auth"
	"user-profiles/internal/http/cookie"
)

func StrictAuthMiddleware(jwtService auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessCookie, err := cookie.Get("__up_access_token", r)
			if accessCookie == "" {
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}

			token := accessCookie
			claims, err := jwtService.Parse(token)
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
			ctx := context.WithValue(r.Context(), utils.UserIdKey, uId)

			role, ok := claims["role"].(string)
			if !ok {
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}
			// send UserRoleKey in context
			ctx = context.WithValue(ctx, utils.UserRoleKey, role)
			ban, ok := claims["banned"].(bool)
			if !ok {
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}
			// send UseBanKey in context
			ctx = context.WithValue(ctx, utils.UserBanKey, ban)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
