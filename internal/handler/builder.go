package handler

import (
	"github.com/Kanbenn/gophermart/internal/app"
	"github.com/Kanbenn/gophermart/internal/config"
	"github.com/Kanbenn/gophermart/internal/storage"
)

type Handler struct {
	cfg config.Config
	app *app.App
	db  *storage.Pg
}

func New(cfg config.Config, app *app.App, db *storage.Pg) *Handler {
	return &Handler{cfg, app, db}
}
