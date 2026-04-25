package middleware

import (
	"net/http"
	"strconv"
	"user-profiles/internal/auth"
)

func CheckAuthAndRedirect(jwtService auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uId, err := GetUserID(r.Context())
			if err == nil {
				http.Redirect(w, r, "/profile/"+strconv.Itoa(uId), http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
