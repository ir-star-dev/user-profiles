package middlewares

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func CSRF(isProd bool) func(http.Handler) http.Handler {
	return csrf.Protect(
		[]byte("32-byte-long-super-secret-key!"),

		// true Prod
		csrf.Secure(isProd),
		csrf.TrustedOrigins([]string{
			"localhost:8080",
			"127.0.0.1:8080",
		}),

		// Cookie settings
		csrf.CookieName("csrf_token"),
		csrf.Path("/"),
		csrf.HttpOnly(true),
		csrf.SameSite(csrf.SameSiteLaxMode),
	)
}
