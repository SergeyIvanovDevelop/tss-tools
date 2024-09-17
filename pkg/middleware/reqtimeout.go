package middleware

import (
	"context"
	"net/http"
	"time"
)

const timeout = 5 * time.Minute

func requestTimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Создаем контекст с таймаутом
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		// Создаем новый запрос с обновленным контекстом
		r = r.WithContext(ctx)

		// Завершаем выполнение запроса, если произошел таймаут
		done := make(chan struct{})
		go func() {
			defer close(done)
			next.ServeHTTP(w, r)
		}()

		select {
		case <-ctx.Done():
			// Обработка таймаута (запрос был отменен)
			w.WriteHeader(http.StatusGatewayTimeout)
			w.Write([]byte("Request timed out"))
		case <-done:
			// Завершение обработки запроса
		}
	})
}
