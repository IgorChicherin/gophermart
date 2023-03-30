package db

import (
	"database/sql"
	"embed"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"net/http"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(databaseDSN string) error {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return err
	}
	driver, err := pgx.WithInstance(db, &pgx.Config{})

	if err != nil {
		return err
	}

	source, err := httpfs.New(http.FS(migrations), "migrations")

	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(
		"migrations", source, "postgres", driver)

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
