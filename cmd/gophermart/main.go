package main

import (
	"log"
	"net/http"

	"github.com/Kanbenn/gophermart/internal/config"
	"github.com/Kanbenn/gophermart/internal/handler"
	"github.com/Kanbenn/gophermart/internal/router"
	"github.com/Kanbenn/gophermart/internal/storage"
)

func main() {

	cfg := config.NewFromFlagsAndEnvs()

	pg := storage.NewPostgres(cfg)
	defer pg.Close()

	h := handler.New(cfg, pg)
	r := router.New(h)

	log.Println("starting web-server on address:", cfg.Addr)
	err := http.ListenAndServe(cfg.Addr, r)
	if err != nil {
		panic(err)
	}
}
