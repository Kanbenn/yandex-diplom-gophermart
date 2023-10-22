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

	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	// defer stop()
	go pg.LaunchWorkerAccrual()

	h := handler.New(cfg, pg)
	r := router.New(h)

	log.Println("starting web-server on address:", cfg.Addr)
	err := http.ListenAndServe(cfg.Addr, r)
	if err != nil {
		panic(err)
	}
	// <-ctx.Done()
	// pg.StopCh <- struct{}{}
	// os.Exit(0)

	// c :=
	// go func() {
	// 	fmt.Println("got the signal from OS", <-c)
	// 	pg.StopCh<- struct{}{}
	// 	fmt.Println("Exiting...")
	// 	os.Exit(0)
	// }()
}

// TODO:
// graceful shutdown с каналами и мидлваркой
// при запуске, подгружать необработанные заказы из базы в горутину для accrual.
