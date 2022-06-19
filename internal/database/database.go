package database

import (
	"database/sql"
	"flag"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	databaseURL     = flag.String("database-url", "postgres://postgres:password@db:5432/postgres?sslmode=disable", "Database URL")
	databaseTimeout = flag.Int64("database-timeout-ms", 2000, "")
)

func Connect() (*sqlx.DB, error) {
	dbURL := *databaseURL

	logrus.Debug("Connecting to Database")
	conn, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to database")
	}
	conn.SetMaxOpenConns(32)

	if err := waitForDB(conn.DB); err != nil {
		return nil, err
	}

	if err := migrateDB(conn.DB); err != nil {
		return nil, errors.Wrap(err, "Could not migrate to database")
	}

	return conn, nil
}

func New() (Database, error) {
	conn, err := Connect()
	if err != nil {
		return nil, err
	}

	d := &database{
		conn: conn,
	}
	return d, nil
}

func waitForDB(conn *sql.DB) error {
	ready := make(chan struct{})

	go func() {
		for {
			if err := conn.Ping(); err == nil {
				close(ready)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	select {
	case <-ready:
		return nil
	case <-time.After(time.Duration(*databaseTimeout) * time.Millisecond):
		return errors.New("database not ready")
	}
}
