package worker

import (
	"time"

	"github.com/Kanbenn/gophermart/internal/config"
	"github.com/Kanbenn/gophermart/internal/models"
)

type Worker struct {
	cfg config.Config
	s   storer
}

type storer interface {
	SelectOrdersForAccrual() (orders []models.Accrual, err error)
	UpdateOrderStatusAndUserBalance(order models.Accrual)
}

type result struct {
	statusCode int
	delay      time.Duration
}

func New(cfg config.Config, s storer) *Worker {
	return &Worker{cfg, s}
}
