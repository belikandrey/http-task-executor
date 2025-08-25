package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"task-service/internal/logger"
	"time"
)

func New(log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()

			defer func() {
				log.Infof("%s request to %s completed with status code %d, bytes = %d, duration = %s. Remote Addr : %s, User Agent : %s, Request id %s",
					r.Method,
					r.URL.Path,
					ww.Status(),
					ww.BytesWritten(),
					time.Since(start).String(),
					r.RemoteAddr,
					r.UserAgent(),
					middleware.GetReqID(r.Context()))
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
