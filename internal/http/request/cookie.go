package request

import (
	"net/http"
	"time"
)

func GetCookie(name string, r *http.Request) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func SetCookie(token string, name string, ttl time.Duration, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    token,
		HttpOnly: true,
		Secure:   false, // true
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(ttl),
	})
}
