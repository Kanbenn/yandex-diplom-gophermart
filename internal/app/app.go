package app

import (
	"log"

	"github.com/Kanbenn/gophermart/internal/models"
	"github.com/Kanbenn/gophermart/pkg/luhn"
)

func (app *App) UserRegister(user models.User) (uid int, err error) {
	return app.s.InsertNewUser(user)
}
func (app *App) UserAuth(in models.User) (uid int, err error) {
	out, err := app.s.SelectUserAuth(in.Login)
	if err != nil {
		log.Printf("app.UserAuth unexpected error: %#v \n", err)
		return out.ID, models.ErrUnxpectedError
	}
	if in.Password != out.Password {
		log.Println("app.UserAuth error: wrong login or password for", in.Login)
		return out.ID, models.ErrUserUnknown
	}
	return out.ID, nil
}
func (app *App) OrderNew(order models.Order) (err error) {
	if !luhn.IsValidLuhnNumber([]byte(order.Number)) {
		log.Println("app.OrderNew luhn formula error:", order)
		return models.ErrLuhnFormulaViolation
	}

	if order.IsWithdrawal {
		err = app.s.InsertOrderWithBonus(order)
		log.Println("app.OrderNew with Bonus", order, err)
		return err
	}

	log.Println("app.OrderNew to Accrual", order, err)
	err = app.s.InsertOrder(order)
	if err != nil {
		return err
	}
	app.wa.NotifyWorker(order)
	return nil
}
func (app *App) UserBalance(uid int) (ub models.UserBalance, err error) {
	return app.s.SelectUserBalance(uid)
}
func (app *App) UserOrders(uid int) (orders []models.OrderResponse, err error) {
	orders, err = app.s.SelectUserOrders(uid)
	if err != nil {
		log.Println("app.UserOrders unexpected error:", orders, err)
		return orders, models.ErrUnxpectedError
	}
	if len(orders) < 1 {
		return orders, models.ErrNoContent
	}
	return orders, nil
}
func (app *App) UserWithdrawHistory(uid int) (orders []models.Order, err error) {
	orders, err = app.s.SelectUserWithdrawHistory(uid)
	if err != nil {
		log.Println("app.UserWithdrawHistory err:", orders, err)
		return orders, models.ErrUnxpectedError
	}
	if len(orders) < 1 {
		log.Println("h.GetUserHistory err: нет ни одного списания", uid, orders, err)
		return orders, models.ErrNoContent
	}
	return orders, nil
}
