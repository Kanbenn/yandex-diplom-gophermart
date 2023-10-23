package repo

// import (
// 	"github.com/Kanbenn/gophermart/internal/config"
// 	"github.com/Kanbenn/gophermart/internal/models"
// )

// type Repo struct {
// 	Cfg    config.Config
// 	StopCh chan struct{}
// 	s      storer
// }

// type storer interface {
// 	InsertOrderWithBonus(o models.OrderNew) error
// 	InsertOrder(o models.OrderNew) error
// 	InsertUser(user models.UserInsert) (uid int, err error)
// 	SelectUserAuth(login string) (user models.UserInsert)
// 	SelectUserBalance(uid int) (ub models.UserBalance, err error)
// 	SelectUserWithdrawHistory(uid int) (orders []models.OrderNew, err error)
// 	SelectUserAllOrders(uid int) (orders []models.OrderResponse, err error)
// }

// func New(cfg config.Config, s storer) *Repo {
// 	// repo :=
// 	sch := make(chan struct{})
// 	return &Repo{cfg, sch, s}
// }

// func (r *Repo) LaunchWorkerAccrual() {

// }
