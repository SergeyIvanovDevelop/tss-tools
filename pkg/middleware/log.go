package middleware

import (
	"net/http"
	"time"

	log "github.com/SergeyIvanovDevelop/tss-tools/pkg/logger"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		// Запоминаем первичные настройки логгера
		oldLogger := log.Logger

		// Добавляем поле с "request_id" в логгер для данного запроса
		requestID, ok := r.Context().Value(requestIDKey).(string)
		if ok {
			log.Logger = log.WithFields(log.Fields{
				"request_id": requestID,
			})
		}

		next.ServeHTTP(&lw, r)

		uri := r.RequestURI
		method := r.Method

		duration := time.Since(start)

		log.WithFields(log.Fields{
			"uri":         uri,
			"method":      method,
			"duration":    duration,
			"status code": responseData.status,
			"answer size": responseData.size,
		}).Infof("logMiddleware")

		// Возвращаем первичные настройки логгера
		log.Logger = oldLogger
	})
}
