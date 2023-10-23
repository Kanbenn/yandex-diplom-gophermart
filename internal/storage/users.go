package storage

import (
	"github.com/Kanbenn/gophermart/internal/models"
)

func (pg *Pg) InsertNewUser(user models.UserInsert) (uid int, err error) {
	q := `INSERT INTO users(login, password) VALUES($1, $2)
		  RETURNING id`
	err = pg.Sqlx.QueryRowx(q, user.Login, user.Password).Scan(&uid)
	if isNotUniqueInsert(err) {
		return uid, ErrLoginNotUnique
	}
	return uid, err
}

func (pg *Pg) SelectUserAuth(login string) (user models.UserInsert) {
	q := `SELECT id, login, password FROM users WHERE login = $1`
	_ = pg.Sqlx.Get(&user, q, login)
	return user
}

func (pg *Pg) SelectUserAllOrders(uid int) (orders []models.OrderResponse, err error) {
	q := `
	SELECT number,status,bonus,time FROM orders
	WHERE user_id = $1 ORDER BY created_at`
	err = pg.Sqlx.Select(&orders, q, uid)
	return orders, err
}

func (pg *Pg) SelectUserBalance(uid int) (ub models.UserBalance, err error) {
	q := `SELECT id, balance, withdrawn FROM users WHERE id = $1`
	err = pg.Sqlx.Get(&ub, q, uid)
	return ub, err
}

func (pg *Pg) SelectUserWithdrawHistory(uid int) (orders []models.OrderNew, err error) {
	q := `SELECT number, sum, time FROM orders WHERE user_id = $1 AND sum > 0`
	err = pg.Sqlx.Select(&orders, q, uid)
	return orders, err
}
