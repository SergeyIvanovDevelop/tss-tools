package migrations

import (

	// need for migrations
	"github.com/golang-migrate/migrate/v4"
	// need for migrations
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// need for migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"
	// need for migrations
	_ "github.com/lib/pq"

	log "github.com/SergeyIvanovDevelop/tss-tools/pkg/logger"
)

var pkgLog log.Log

const pkgName = "tss-tools/pkg/migrations"

func ApplyMigrations(connString, pathToMigrations string) error {
	fncLogger := log.AddLoggerFields(pkgLog, pkgName, log.Fields{
		"func": "ApplyMigrations",
	})

	m, err := migrate.New(
		pathToMigrations,
		connString,
	)
	if err != nil {
		return fncLogger.WrapError("ошибка инициализации миграций connection string '%v', path to migrations '%s': %w", connString, pathToMigrations, err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fncLogger.WrapError("ошибка применения миграций connection string '%v', path to migrations '%s': %w", connString, pathToMigrations, err)
	}

	fncLogger.Info("Migrations applied successfully!")

	return nil

}
