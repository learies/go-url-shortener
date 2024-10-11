package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/learies/go-url-shortener/internal/logger"
)

type ResponseWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData  *responseData
		headerWritten bool
		mu            sync.Mutex
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size

	r.mu.Lock()
	if !r.headerWritten {
		r.responseData.status = http.StatusOK
		r.headerWritten = true
	}
	r.mu.Unlock()

	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.headerWritten {
		return
	}
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
	r.headerWritten = true
}

func WithLogging(next http.Handler) http.Handler {
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
			"duration", duration.Milliseconds(),
			"size", responseData.size,
		)
	})
}
