package cookie

import (
	"net/http"
	"time"
)

func Get(name string, r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, err
		}
		return nil, err
	}
	return cookie, nil
}

func Set(token string, name string, ttl time.Duration, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(ttl),
	})
}
