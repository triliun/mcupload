package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/triliun/mcupload/backend/logger"
	"go.uber.org/zap"
)

// Logger returns a middleware for logging
func Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response for capture status code and body size
		wrappedWriter := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		var requestBody []byte
		if r.Body != nil {
			requestBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Process request
		next.ServeHTTP(wrappedWriter, r)

		latency := time.Since(start)
		statusCode := wrappedWriter.Status()

		// Log request
		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
			zap.String("ip", r.RemoteAddr), // from middleware.RealIP
			zap.String("user_agent", r.UserAgent()),
			zap.Duration("latency", latency),
			zap.String("origin", r.Header.Get("Origin")),
			zap.String("request_id", middleware.GetReqID(r.Context())),
			zap.Int("response_size", wrappedWriter.BytesWritten()),
			zap.Int("request_size", len(requestBody)),
		}

		switch {
		case statusCode >= 500:
			logger.Error("Server error", fields...)
		case statusCode >= 400:
			logger.Error("Client error", fields...)
		default:
			logger.Info("Request completed", fields...)
		}
	}

	return http.HandlerFunc(fn)
}
