package requestlogger

import (
	"net/http"
	"time"

	"github.com/Shoyeb45/server/pkg/logger"
	"github.com/go-chi/chi/middleware"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		logger.Log.Info("request",
			"request_id", middleware.GetReqID(r.Context()),
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"bytes", ww.BytesWritten(),
			"latency", time.Since(start).String(),
			"ip", r.RemoteAddr,
		)
	})
}
