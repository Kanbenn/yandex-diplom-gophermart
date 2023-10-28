package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/Kanbenn/gophermart/internal/app"
	"github.com/Kanbenn/gophermart/internal/config"
	"github.com/Kanbenn/gophermart/internal/postgres"
	"github.com/Kanbenn/gophermart/internal/server"
	"github.com/Kanbenn/gophermart/internal/worker"
)

func main() {

	cfg := config.NewFromFlagsAndEnvs()

	pg := postgres.New(cfg)
	defer pg.Close()

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	wa := worker.New(cfg, pg)
	go wa.LaunchWorkerAccrual(ctx)

	app := app.New(cfg, pg, wa)

	srv := server.New(cfg, app)
	go srv.ShutdownOnSignal(ctx)
	srv.Launch()

}
