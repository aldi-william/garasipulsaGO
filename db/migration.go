package db

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
)

const (
	migrationSource = "file://./db/migrations"
)

func autoMigrate(db *sql.DB) error {
	// driver, err := postgres.WithInstance(db, &postgres.Config{})
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return errors.Wrap(err, "unable to create migration driver")
	}
	mig, err := migrate.NewWithDatabaseInstance(migrationSource, driverName, driver)
	if err != nil {
		return errors.Wrap(err, "unable to init migration")
	}
	err = mig.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "failed migration up")
	}
	return nil
}
