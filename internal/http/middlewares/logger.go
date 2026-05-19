package middlewares

import (
	"log/slog"
	"net/http"
	//"runtime"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func Logger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{
				ResponseWriter: w,
				status:         200,
			}
			defer func() {
				duration := time.Since(start)
				log.Info("http request",
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
					slog.Int("status", rec.status),
					slog.Duration("duration", duration),
				)
			}()
			next.ServeHTTP(rec, r)
		})
	}
}