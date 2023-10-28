package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Kanbenn/gophermart/internal/app"
	"github.com/Kanbenn/gophermart/internal/config"
	"github.com/Kanbenn/gophermart/internal/handler"
	"github.com/Kanbenn/gophermart/internal/router"
)

type Server struct {
	http.Server
}

func New(cfg config.Config, app *app.App) *Server {
	h := handler.New(cfg, app)
	r := router.New(h)

	srv := Server{
		http.Server{
			Addr:    cfg.Addr,
			Handler: r}}
	return &srv
}

func (srv *Server) ShutdownOnSignal(ctxFromSignal context.Context) {
	<-ctxFromSignal.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	log.Println("shutting down server..")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("failed to shutdown the server gracefully, forcing exit", err)
	}
}

func (srv *Server) Launch() {
	log.Println("starting web-server on address:", srv.Server.Addr)
	if err := srv.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("failed to start listening", err)
	}
}
