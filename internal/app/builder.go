package app

import (
	"github.com/Kanbenn/gophermart/internal/config"
	"github.com/Kanbenn/gophermart/internal/models"
)

type storer interface {
	InsertNewUser(user models.User) (uid int, err error)
	InsertOrder(o models.Order) error
	InsertOrderWithBonus(o models.Order) error
	SelectUserAuth(login string) (user models.User, err error)
	SelectUserBalance(uid int) (ub models.UserBalance, err error)
	SelectUserOrders(uid int) (orders []models.OrderResponse, err error)
	SelectUserWithdrawHistory(uid int) (orders []models.Order, err error)
}

type App struct {
	Cfg    config.Config
	StopCh chan struct{}
	s      storer
}

func New(cfg config.Config, s storer) *App {
	sch := make(chan struct{})
	return &App{cfg, sch, s}
}
