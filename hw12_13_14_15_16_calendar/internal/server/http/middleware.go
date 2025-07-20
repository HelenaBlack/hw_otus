package internalhttp

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func loggingMiddleware(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{w, http.StatusOK}
			next.ServeHTTP(rw, r)

			latency := time.Since(start)

			ip := r.RemoteAddr
			if ipHeader := r.Header.Get("X-Real-IP"); ipHeader != "" {
				ip = ipHeader
			} else if ipHeader := r.Header.Get("X-Forwarded-For"); ipHeader != "" {
				ip = strings.Split(ipHeader, ",")[0]
			}

			timestamp := time.Now().Format("02/Jan/2006:15:04:05 -0700")

			ua := r.UserAgent()

			// Формируем строку лога в формате Apache/Nginx
			logLine := fmt.Sprintf("%s [%s] %s %s %s %d %d \"%s\"",
				ip,
				timestamp,
				r.Method,
				r.RequestURI,
				r.Proto,
				rw.statusCode,
				latency.Milliseconds(),
				ua,
			)
			logger.Info(logLine)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int // сохраняем статус-код для логирования
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
