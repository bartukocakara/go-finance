package database

import (
	"database/sql"
	"finance-app/internal/config"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func migrateDB(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "Connection To Database... log from from Migration module")
	}

	migrationSource := fmt.Sprintf(
		"file://%s/internal/database/migrations/", *config.DataDirectory)
	migrator, err := migrate.NewWithDatabaseInstance(
		migrationSource,
		"postgres", driver)
	if err != nil {
		return errors.Wrap(err, "Creating Migrator")
	}
	logrus.Debug("Passed migration source ", migrationSource)
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "Executing Migration")
	}

	version, dirty, err := migrator.Version()
	if err != nil {
		return errors.Wrap(err, "Getting migration version")
	}

	logrus.WithFields(logrus.Fields{
		"version": version,
		"dirty":   dirty,
	}).Debug("Database migrated")
	return nil
}
