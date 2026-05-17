package middlewares

import (
	"context"
	"errors"
	"net/http"
	"time"
	"user-profiles/internal/auth"
	"user-profiles/internal/http/request"
)

type authData struct {
	userId int
	role   string
	banned bool
}

func StrictAuthMiddleware(jwtService auth.JWTService, authService auth.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authData, err := AuthenticateRequest(w, r, jwtService, authService)

			if err != nil {
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), request.UserIdKey, authData.userId)
			ctx = context.WithValue(ctx, request.UserRoleKey, authData.role)
			ctx = context.WithValue(ctx, request.UserBanKey, authData.banned)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func SoftAuthMiddleware(jwtService auth.JWTService, authService auth.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authData, err := AuthenticateRequest(w, r, jwtService, authService)
			// guest allowed
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), request.UserIdKey, authData.userId)
			ctx = context.WithValue(ctx, request.UserRoleKey, authData.role)
			ctx = context.WithValue(ctx, request.UserBanKey, authData.banned)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthenticateRequest(w http.ResponseWriter, r *http.Request, jwtService auth.JWTService, authService auth.AuthService) (*authData, error) {
	var accessToken string

	accessCookie, _ := request.GetCookie("__up_access_token", r)
	refreshCookie, _ := request.GetCookie("__up_refresh_token", r)

	// access exists
	if accessCookie != "" {
		accessToken = accessCookie
	}

	// try refresh
	if accessToken == "" && refreshCookie != "" {
		tokens, _ := authService.Refresh(refreshCookie)

		if tokens != nil {
			accessToken = tokens.Access

			request.SetCookie(tokens.Access, "__up_access_token", 5*time.Minute, w)
			request.SetCookie(tokens.Refresh, "__up_refresh_token", 7*24*time.Hour, w)
		} else {
			request.SetCookie("", "__up_access_token", -time.Minute, w)
			request.SetCookie("", "__up_refresh_token", -time.Minute, w)
		}
	}

	if accessToken == "" {
		return nil, errors.New("Unauthorized")
	}

	claims, err := jwtService.Parse(accessToken)
	if err != nil {
		return nil, err
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return nil, errors.New("Invalid sub")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("Invalid role")
	}

	ban, ok := claims["banned"].(bool)
	if !ok {
		return nil, errors.New("Invalid banned")
	}

	return &authData{
		userId: int(sub),
		role:   role,
		banned: ban}, nil
}
