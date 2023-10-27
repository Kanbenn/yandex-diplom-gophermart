package main

import (
	"log"
	"net/http"

	"github.com/Kanbenn/gophermart/internal/app"
	"github.com/Kanbenn/gophermart/internal/config"
	"github.com/Kanbenn/gophermart/internal/handler"
	"github.com/Kanbenn/gophermart/internal/router"
	"github.com/Kanbenn/gophermart/internal/storage"
	"github.com/jmoiron/sqlx"
)

func main() {

	cfg := config.NewFromFlagsAndEnvs()
	conn := connectToPostgres(cfg)
	pg := storage.NewInPostgres(cfg, conn)
	defer pg.Close()
	pg.CreateTables()

	go pg.LaunchWorkerAccrual()

	app := app.New(cfg, pg)

	h := handler.New(cfg, app, pg)
	r := router.New(h)

	log.Println("starting web-server on address:", cfg.Addr)
	err := http.ListenAndServe(cfg.Addr, r)
	if err != nil {
		panic(err)
	}

}

func connectToPostgres(cfg config.Config) *sqlx.DB {

	conn, err := sqlx.Open("postgres", cfg.PgConnStr)
	if err != nil {
		log.Fatal("error at connecting to Postgres:", cfg.PgConnStr, err)
	}
	if err := conn.Ping(); err != nil {
		log.Fatal("error at pinging the Postgres db:", cfg.PgConnStr, err)
	}
	return conn
}
