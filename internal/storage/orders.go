package storage

import (
	"log"

	"github.com/Kanbenn/gophermart/internal/models"
)

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
	return checkInsertOrderResults(r)
}

func (pg *Pg) InsertOrderWithBonus(o models.Order) error {
	tx, err := pg.Sqlx.Beginx()
	if err != nil {
		log.Printf("\n pg.InsertOrderWithBonus unexpected tx error: %#v \n\n", err)
		return models.ErrUnxpectedError
	}
	defer tx.Rollback()

	qo := `
	INSERT INTO orders (number, user_id, status, sum, time) 
	VALUES(:number, :user, :status, :sum, :time)`
	o.Time = pg.timeForNewOrders()
	o.Status = pg.statusForNewOrderWithBonus()
	if _, err := tx.NamedExec(qo, o); err != nil {
		return checkInserOrderWithBonusErr(err)
	}

	qb := `
	UPDATE users SET balance = users.balance - :sum, withdrawn = users.withdrawn + :sum 
	WHERE id=:user`
	_, err = tx.NamedExec(qb, o)
	if err != nil {
		return checkUpdateUserBalanceErr(err)
	}

	tx.Commit()
	log.Println("успешно принят новый бонусный заказ", o)
	return nil
}
