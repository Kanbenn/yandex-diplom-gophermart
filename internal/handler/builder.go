package handler

import (
	"github.com/Kanbenn/gophermart/internal/config"
	"github.com/Kanbenn/gophermart/internal/storage"
)

type Handler struct {
	cfg config.Config
	db  *storage.Pg
}

func New(cfg config.Config, db *storage.Pg) *Handler {
	return &Handler{cfg, db}
}
