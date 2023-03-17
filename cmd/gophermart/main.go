package main

import (
	"context"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/config"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/router"
	"github.com/IgorChicherin/gophermart/internal/pkg/authlib/sha256"
	"github.com/IgorChicherin/gophermart/internal/pkg/db"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.GetSeverConfig()

	if cfg.IsDefaultHashKey() {
		log.Warning("default secret key has been used")
	}

	ctxDB := context.Background()
	conn, err := pgx.Connect(ctxDB, cfg.DatabaseURI)

	if err != nil {
		log.Fatalln(err)
	}

	if err := db.Migrate(cfg.DatabaseURI); err != nil {
		log.Fatalln(err)
	}
	defer conn.Close(ctxDB)

	err = conn.Ping(ctxDB)

	if err != nil {
		log.Fatalln(err)
	}

	hashService := sha256.NewSha256HashService(cfg.HashKey)

	r := router.NewRouter(conn, hashService)

	if err := r.Run(cfg.Address); err != nil {
		log.Fatalln(err)
	}
}
