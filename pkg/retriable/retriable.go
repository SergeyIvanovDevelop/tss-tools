package retriable

import (
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"

	log "github.com/SergeyIvanovDevelop/tss-tools/pkg/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var pkgLog log.Log

const pkgName = "pkg/retriable"

// retryableOperation представляет собой функцию, которую нужно повторить
type retryableOperation func() error

// retryWithBackoff выполняет операцию с повторными попытками в случае retriable ошибок
func RetryWithBackoff(operation retryableOperation, retryNumber int) error {
	var fncLogger = log.AddLoggerFields(pkgLog, pkgName, log.Fields{
		"func": "RetryWithBackoff",
	})
	const step = 2
	var maxAttempts = retryNumber + 1
	var lastErr error

	var backoffIntervals []time.Duration = make([]time.Duration, 0, retryNumber)

	// Интервалы между попытками
	backoffIntervals = append(backoffIntervals, 1*time.Second)
	for i := 1; i < retryNumber; i++ {
		interval := time.Duration(i+step) * time.Second
		backoffIntervals = append(backoffIntervals, interval)
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Выполняем операцию
		lastErr = operation()

		if lastErr == nil {
			return nil // Успешное выполнение без ошибок
		}

		// Проверяем, является ли ошибка retryable
		if isRetryableError(lastErr) {
			if attempt < maxAttempts {
				fncLogger.Errorf("Попытка %d не удалась, ожидаем %v перед повторной попыткой...\n", attempt, backoffIntervals[attempt-1])
				time.Sleep(backoffIntervals[attempt-1])
			}
		} else {
			return lastErr // Если ошибка не подлежит повторной попытке, возвращаем её
		}
	}

	return fmt.Errorf("операция не удалась после %d попыток: %w", maxAttempts, lastErr)
}

// isRetryableError определяет, является ли ошибка повторяемой
func isRetryableError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// Проверяем на конкретный код ошибки PostgreSQL
		switch pgErr.Code {
		case pgerrcode.ConnectionException, // ошибка соединения
			pgerrcode.ConnectionFailure,
			pgerrcode.ConnectionDoesNotExist,
			pgerrcode.SQLClientUnableToEstablishSQLConnection,
			pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection,
			pgerrcode.TransactionResolutionUnknown,
			pgerrcode.ProtocolViolation,
			pgerrcode.UniqueViolation,
			pgerrcode.AdminShutdown: // сбой в работе сервера
			return true
		}
	}

	// Проверка на сетевые ошибки
	var netErr net.Error // указатель не нужен, т.к. net.Error - уже является интерфейсом (добавление указателя - избыточно)
	if errors.As(err, &netErr) {
		return true
	}

	// Проверка на ошибки доступа к файлам
	if errors.Is(err, os.ErrPermission) {
		return true // ошибка разрешения доступа к файлу
	}

	// Проверка на блокировку файла
	var errno *syscall.Errno
	if errors.As(err, &errno) {
		if *errno == syscall.EAGAIN {
			return true // ошибка блокировки файла
		}
	}

	// Добавляем другие пользовательские условия для повторяемых ошибок
	return false
}
