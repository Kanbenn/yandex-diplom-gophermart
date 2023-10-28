package handler

import (
	"github.com/Kanbenn/gophermart/internal/app"
	"github.com/Kanbenn/gophermart/internal/config"
)

type Handler struct {
	cfg config.Config
	app *app.App
}

func New(cfg config.Config, app *app.App) *Handler {
	return &Handler{cfg, app}
}
