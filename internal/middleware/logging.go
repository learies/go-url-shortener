package middleware

import (
	"net/http"
	"time"

	"github.com/learies/go-url-shortener/internal/config/logger"
)

// responseData encapsulates HTTP response details.
type responseData struct {
	status int
	size   int
}

// loggingResponseWriter is a custom response writer that captures HTTP response details.
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData  *responseData
	headerWritten bool
}

// Header returns the header map that will be sent by WriteHeader.
func (r *loggingResponseWriter) Header() http.Header {
	return r.ResponseWriter.Header()
}

// Write captures the size of the response and writes the response.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	if err != nil {
		return size, err
	}
	r.responseData.size += size
	// To ensure that the status code is captured correctly, we only set it if it hasn't been set before.
	if !r.headerWritten {
		r.responseData.status = http.StatusOK
		r.headerWritten = true
	}
	return size, err
}

// WriteHeader captures the status code and writes the header.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// To ensure that the status code is captured correctly, we only set it if it hasn't been set before.
	if !r.headerWritten {
		r.ResponseWriter.WriteHeader(statusCode)
		r.responseData.status = statusCode
		r.headerWritten = true
	}
}

// MiddlewareLogger logs details about each HTTP request.
func MiddlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := &loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
			headerWritten:  false,
		}
		next.ServeHTTP(lw, r)
		duration := time.Since(start)

		logger.Log.Info("Request completed",
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	})
}
