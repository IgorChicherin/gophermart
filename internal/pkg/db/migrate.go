package db

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(databaseDSN string) error {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return err
	}
	driver, err := pgx.WithInstance(db, &pgx.Config{})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/pkg/db/migrations",
		"postgres", driver)

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
