package authserv

import (
	"context"
	"net/http"
	"time"

	"github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/handlers"
	"github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/repository"
	log "github.com/SergeyIvanovDevelop/tss-tools/pkg/logger"

	"github.com/gorilla/mux"
)

var pkgLog log.Log

const pkgName string = "tss-tools/pkg/authserv"

// ServerConfig содержит параметры для настройки сервера.
type ServerConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Run запускает HTTP сервер в отдельной горутине с поддержкой graceful-shutdown.
func Run(ctx context.Context, db repository.AuthRepository, config ServerConfig) error {
	fncLogger := log.AddLoggerFields(pkgLog, pkgName, log.Fields{
		"func": "Run",
	})

	r := mux.NewRouter()

	r.HandleFunc("/api/user/register", handlers.Register(db)).Methods("POST")
	r.HandleFunc("/api/user/login", handlers.Login(db)).Methods("POST")
	r.HandleFunc("/api/user/revoke", handlers.Revoke(db)).Methods("POST")
	r.HandleFunc("/api/user/validate", handlers.Validate(db)).Methods("POST")

	srv := &http.Server{
		Addr:         config.Addr,
		Handler:      r,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	startBlacklistCleaner(db)

	serverErrors := make(chan error, 1)
	go func() {
		fncLogger.Infof("Сервер аутентификации запущен на %s", config.Addr)
		serverErrors <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		fncLogger.Info("Получен сигнал завершения работы, выключаем сервер аутентификации...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			fncLogger.Error("Ошибка при завершении работы сервера аутентификации: %v", err)
			return err
		}
		fncLogger.Info("Сервер аутентификации успешно завершил работу")
		return nil
	case err := <-serverErrors:
		return err
	}
}

func startBlacklistCleaner(repo repository.AuthRepository) {
	fncLogger := log.AddLoggerFields(pkgLog, pkgName, log.Fields{
		"func": "startBlacklistCleaner",
	})
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := repo.CleanExpiredTokens(); err != nil {
					fncLogger.Errorf("Ошибка очистки черного списка: %v", err)
				}
			}
		}
	}()
}
