package main

import (
	"context"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/router"
	"github.com/IgorChicherin/gophermart/internal/pkg/authlib/sha256"
	"github.com/IgorChicherin/gophermart/internal/pkg/db"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

func main() {

	databaseDSN := "postgres://test:test@localhost:5432/gophermart?sslmode=disable"
	ctxDB := context.Background()

	conn, err := pgx.Connect(ctxDB, databaseDSN)

	if err != nil {
		log.Fatalln(err)
	}

	if err := db.Migrate(databaseDSN); err != nil {
		log.Fatalln(err)
	}
	defer conn.Close(ctxDB)

	err = conn.Ping(ctxDB)

	if err != nil {
		log.Fatalln(err)
	}

	hashService := sha256.NewSha256HashService("12313")

	r := router.NewRouter(conn, hashService)

	if err := r.Run("localhost:8080"); err != nil {
		log.Fatalln(err)
	}
}
