package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"http-task-executor/task-service/internal/task-service/logger"
)

// New creates new logger middleware that adds logging in http requests.
func New(log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()

			defer func() {
				log.Infof(
					"%s request to %s completed with status code %d, bytes = %d, duration = %s. Remote Addr : %s, User Agent : %s, Request id %s",
					r.Method,
					r.URL.Path,
					ww.Status(),
					ww.BytesWritten(),
					time.Since(start).String(),
					r.RemoteAddr,
					r.UserAgent(),
					middleware.GetReqID(r.Context()),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
