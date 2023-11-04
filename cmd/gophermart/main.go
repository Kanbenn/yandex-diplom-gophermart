package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kanbenn/gophermart/internal/app"
	"github.com/Kanbenn/gophermart/internal/config"
	"github.com/Kanbenn/gophermart/internal/postgres"
	"github.com/Kanbenn/gophermart/internal/server"
	"github.com/Kanbenn/gophermart/internal/worker"
)

func main() {

	cfg := config.New()
	cfg.ParseFlagsAndEnvs()

	pg := postgres.New(cfg)
	defer pg.Close()

	app := app.New(cfg, pg)

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	wa := worker.New(cfg, pg)
	go wa.LaunchWorkerAccrual(ctx)

	srv := server.New(cfg, app)
	go srv.ShutdownOnSignal(ctx)
	srv.Launch()
}
