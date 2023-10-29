package postgres

import (
	"fmt"
	"log"
	"strings"

	"github.com/Kanbenn/gophermart/internal/models"
)

func (pg *Pg) InsertNewUser(user models.User) (uid int, err error) {
	q := `INSERT INTO users(login, password) VALUES($1, $2)
		  RETURNING id`
	err = pg.Sqlx.QueryRowx(q, user.Login, user.Password).Scan(&uid)
	if isNotUniqueInsert(err) {
		return uid, models.ErrLoginNotUnique
	}
	return uid, err
}

func (pg *Pg) SelectUserAuth(login string) (user models.User, err error) {
	q := `SELECT id, login, password FROM users WHERE login = $1`
	err = pg.Sqlx.Get(&user, q, login)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (pg *Pg) SelectUserOrders(uid int) (orders []models.UserOrder, err error) {
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

func (pg *Pg) SelectUserWithdrawHistory(uid int) (orders []models.Order, err error) {
	q := `SELECT number, sum, time FROM orders WHERE user_id = $1 AND sum > 0`
	err = pg.Sqlx.Select(&orders, q, uid)
	return orders, err
}

func (pg *Pg) InsertOrder(o models.Order) error {
	q := `
	INSERT INTO orders (number, user_id, time) VALUES($1, $2, $3)
	ON CONFLICT (number) DO UPDATE SET updated_at = NOW()
	RETURNING user_id, number, (updated_at > created_at) AS was_conflict`

	o.Time = pg.timeForNewOrders()
	r := orderInsertResult{}
	r.Err = pg.Sqlx.QueryRowx(q, o.Number, o.User, o.Time).
		Scan(&r.User, &r.Number, &r.WasConflict)
	r.SameUser = o.User == r.User
	return pg.checkInsertOrderResults(r)
}

func (pg *Pg) InsertOrderWithdrawal(o models.Order) error {
	tx, err := pg.Sqlx.Beginx()
	if err != nil {
		log.Printf("\n pg.InsertOrderWithdrawal unexpected tx error: %#v \n\n", err)
		return models.ErrUnxpectedError
	}
	defer tx.Rollback()

	qo := `
	INSERT INTO orders (number, user_id, status, sum, time) 
	VALUES(:number, :user, :status, :sum, :time)`
	o.Time = pg.timeForNewOrders()
	o.Status = pg.statusForNewOrderWithBonus()
	if _, err := tx.NamedExec(qo, o); err != nil {
		return pg.checkInserOrderWithBonusErr(err)
	}

	qb := `
	UPDATE users SET balance = users.balance - :sum, withdrawn = users.withdrawn + :sum 
	WHERE id=:user`
	_, err = tx.NamedExec(qb, o)
	if err != nil {
		return pg.checkUpdateUserBalanceErr(err)
	}

	tx.Commit()

	log.Println("успешно принят новый бонусный заказ", o)
	return nil
}

func (pg *Pg) UpdateOrderStatusAndUserBalance(order models.Accrual) {
	tx, err := pg.Sqlx.Beginx()
	if err != nil {
		log.Println("Pg.UpdateOrder: failed to begin the transaction", err)
	}
	defer tx.Rollback()

	log.Println("Pg updating order status:", order)
	qo := `UPDATE orders SET status=:status, bonus=:bonus WHERE number=:number`
	if _, err := tx.NamedExec(qo, order); err != nil {
		log.Println("Pg error updating order in db:", err)
	}

	if pg.isFinalStatus(order.Status) {
		log.Println("Pg updating user balance:", order)
		qb := `UPDATE users SET balance= users.balance + :bonus WHERE id=:user_id`
		if _, err = tx.NamedExec(qb, order); err != nil {
			log.Println("Pg error updating order in db:", err)
		}
	}
	tx.Commit()
}

func (pg *Pg) SelectOrdersForAccrual() (orders []models.Accrual, err error) {
	q := "SELECT number, user_id, status FROM orders WHERE status NOT IN ('%s')"
	finishedStatuses := strings.Join(pg.Cfg.FinishedOrderStatuses, "','")
	q = fmt.Sprintf(q, finishedStatuses)
	err = pg.Sqlx.Select(&orders, q)
	if err != nil {
		return orders, models.ErrNoOrdersForAccrual
	}
	return orders, nil
}
