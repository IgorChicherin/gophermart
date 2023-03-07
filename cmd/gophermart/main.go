package main

import (
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/router"
	log "github.com/sirupsen/logrus"
)

func main() {
	r := router.NewRouter()
	if err := r.Run("localhost:8080"); err != nil {
		log.Fatalln(err)
	}
}
