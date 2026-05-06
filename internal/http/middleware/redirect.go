package middleware

import (
	"net/http"
	"strconv"
	"user-profiles/cmd/user-profiles/auth"
	"user-profiles/internal/http/cookie"
)

func CheckAuthAndRedirect(jwtService auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken, _ := cookie.Get("__up_access_token", r)
			claims, err := jwtService.Parse(accessToken)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			sub, ok := claims["sub"].(int)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}
			http.Redirect(w, r, "/panel/profile/"+strconv.Itoa(sub), http.StatusSeeOther)
		})
	}
}
