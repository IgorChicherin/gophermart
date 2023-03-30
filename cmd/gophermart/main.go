package main

import (
	"context"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/config"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/router"
	"github.com/IgorChicherin/gophermart/internal/pkg/accrual"
	"github.com/IgorChicherin/gophermart/internal/pkg/authlib/sha256"
	"github.com/IgorChicherin/gophermart/internal/pkg/db"
	"github.com/IgorChicherin/gophermart/internal/pkg/moneylib"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.GetSeverConfig()

	if err != nil {
		log.Fatalf("unable to config server: %s", err)
	}

	if cfg.IsDefaultHashKey() {
		log.Warning("default secret key has been used")
	}

	ctxDB := context.Background()
	conn, err := pgx.Connect(ctxDB, cfg.DatabaseURI)

	if err != nil {
		log.Fatalf("unable to connect DB: %s", err)
	}

	if err := db.Migrate(cfg.DatabaseURI); err != nil {
		log.Fatalf("migration failed: %s", err)
	}

	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			log.WithFields(log.Fields{"func": "main"}).Errorln(err)
		}
	}(conn, ctxDB)

	err = conn.Ping(ctxDB)

	if err != nil {
		log.Fatalf("unable to connect DB: %s", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	hashService := sha256.NewSha256HashService(cfg.HashKey)
	moneyService := moneylib.NewMoneyService(100)
	accrualService := accrual.NewAccrualService(ctxDB, conn, cfg.AccrualAddress, moneyService)

	go accrualService.Run()

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router.NewRouter(conn, hashService, moneyService),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Infoln("Server Started")
	<-done
	log.Infoln("Server Stopped")

	ctx := context.Background()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Infoln("Server Exited Properly")

}
