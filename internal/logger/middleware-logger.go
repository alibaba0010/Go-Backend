package logger

import (
	"net/http"
	"time"

	// "sync"
	"go.uber.org/zap"
)
type loggingResponseWriter struct {
    http.ResponseWriter
    status    int  
    // once      sync.Once
}
func (lrw *loggingResponseWriter) WriteHeader(status int) {
    lrw.status = status
    lrw.ResponseWriter.WriteHeader(status)
}
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		// Wrap the original ResponseWriter so we can capture the status code
		lrw := &loggingResponseWriter{ResponseWriter: writer, status: http.StatusOK}
		// Pass the wrapped writer into the next handler so WriteHeader calls
		// set lrw.status and write to the original underlying writer.
		next.ServeHTTP(lrw, request)
		duration := time.Since(start)

		// use the package-level Log variable directly (same package)
		Log.Info("Incoming request",
			zap.String("method", request.Method),
			zap.String("path", request.URL.Path),
			zap.Int("status", lrw.status),
			zap.Duration("duration", duration),
			zap.String("user-agent", request.UserAgent()),
		)
	})
}