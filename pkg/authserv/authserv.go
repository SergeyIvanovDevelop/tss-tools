package authserv

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/handlers"
	"github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/middleware"
	"github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/repository"

	"github.com/gorilla/mux"
)

// ServerConfig содержит параметры для настройки сервера.
type ServerConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Run запускает HTTP сервер в отдельной горутине с поддержкой graceful-shutdown.
func Run(ctx context.Context, db repository.AuthRepository, config ServerConfig) error {
	// Инициализация роутера
	r := mux.NewRouter()

	// Регистрация маршрутов
	r.HandleFunc("/api/user/register", handlers.Register(db)).Methods("POST")
	r.HandleFunc("/api/user/login", handlers.Login(db)).Methods("POST")
	r.HandleFunc("/api/user/revoke", handlers.Revoke(db)).Methods("POST")
	r.HandleFunc("/api/user/validate", handlers.Validate(db)).Methods("POST")

	// Middleware для проверки токенов
	r.Use(middleware.JWTAuthentication)

	// Настройка сервера
	srv := &http.Server{
		Addr:         config.Addr,
		Handler:      r,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	// Запуск планировщика для очистки черного списка токенов
	startBlacklistCleaner(db)

	// Запуск сервера в отдельной горутине
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Сервер запущен на %s", config.Addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// Ожидание завершения через контекст
	select {
	case <-ctx.Done():
		log.Println("Получен сигнал завершения работы, выключаем сервер...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Ошибка при завершении работы сервера: %v", err)
			return err
		}
		log.Println("Сервер успешно завершил работу")
		return nil
	case err := <-serverErrors:
		return err
	}
}

func startBlacklistCleaner(repo repository.AuthRepository) {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := repo.CleanExpiredTokens(); err != nil {
					log.Printf("Ошибка очистки черного списка: %v", err)
				}
			}
		}
	}()
}
