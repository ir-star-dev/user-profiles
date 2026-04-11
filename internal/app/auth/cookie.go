package auth

import (
	"errors"
	"net/http"
	"time"
)

func getRefreshTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", errors.New(NoRefreshtoken)
		}
		return "", errors.New(ReadingCookie)
	}
	return cookie.Value, nil
}

func setCookie(token string, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
}