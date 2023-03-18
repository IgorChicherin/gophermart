package main

import (
	"context"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/config"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/router"
	"github.com/IgorChicherin/gophermart/internal/pkg/accrual"
	"github.com/IgorChicherin/gophermart/internal/pkg/authlib/sha256"
	"github.com/IgorChicherin/gophermart/internal/pkg/db"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	hashService := sha256.NewSha256HashService(cfg.HashKey)
	accrualService := accrual.NewAccrualService(cfg.AccrualAddress)

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router.NewRouter(conn, hashService, accrualService),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Infoln("Server Started")

	<-done
	log.Infoln("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Infoln("Server Exited Properly")

}
