package middlewares

import (
	"net/http"
	"strconv"
	"user-profiles/cmd/user-profiles/auth"
	"user-profiles/internal/http/cookie"
)

func CheckAuthAndRedirect(jwtService auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessCookie, err := cookie.Get("__up_access_token", r)
			if err != nil || accessCookie == "" {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := jwtService.Parse(accessCookie)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			sub, ok := claims["sub"].(float64)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			http.Redirect(w, r, "/panel/profile/"+strconv.Itoa(int(sub)), http.StatusSeeOther)
		})
	}
}
