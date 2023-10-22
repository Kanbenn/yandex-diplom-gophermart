package storage

import (
	"log"
	"time"

	"github.com/Kanbenn/gophermart/internal/models"
)

func (pg *Pg) InsertNewOrder(o models.OrderNew) error {
	q := `
	INSERT INTO orders (number, user_id, time) VALUES($1, $2, $3)
	ON CONFLICT (number) DO UPDATE SET updated_at = NOW()
	RETURNING user_id, (updated_at > created_at) AS was_conflict`

	o.Time = time.Now().Format(time.RFC3339)
	result := models.OrderInsertResult{}
	result.Err = pg.Sqlx.QueryRowx(q, o.Number, o.User, o.Time).
		Scan(&result.User, &result.WasConflict)
	return pg.checkInsertNewOrderResults(result, o)
}

func (pg *Pg) InsertNewBonusOrder(o models.OrderNew) error {
	tx, err := pg.Sqlx.Beginx()
	if err != nil {
		log.Printf("\n pg.InsertNewBonusOrder unexpected tx error: %#v \n\n", err)
		return ErrUnxpectedError
	}
	defer tx.Rollback()

	qu := `SELECT id, balance FROM users WHERE id = $1`
	user := models.UserBalance{}
	if err := tx.Get(&user, qu, o.User); err != nil {
		log.Println("pg.InsertNewBonusOrder", user, err)
		return ErrUserUnknown
	}
	if user.Balance < o.Sum {
		log.Println("pg.InsertNewBonusOrder", user, o, ErrNotEnoughMinerals)
		return ErrNotEnoughMinerals
	}

	qb := `
	UPDATE users SET balance = users.balance - :sum, withdrawn = users.withdrawn + :sum 
	WHERE id = :user`
	_, err = tx.NamedExec(qb, o)
	if err != nil {
		log.Println("pg.InsertNewBonusOrder error updating user's balance:", o, err)
		return ErrUnxpectedError
	}

	qo := `
	INSERT INTO orders (number, user_id, status, sum, time) VALUES($1, $2, 'PROCESSED', $3, $4)
	ON CONFLICT (number) DO UPDATE SET updated_at = NOW()
	RETURNING user_id, (updated_at > created_at) AS was_conflict`
	o.Time = time.Now().Format(time.RFC3339)
	result := models.OrderInsertResult{}
	result.Err = tx.QueryRowx(qo, o.Number, o.User, o.Sum, o.Time).
		Scan(&result.User, &result.WasConflict)

	tx.Commit()
	return pg.checkInsertNewOrderResults(result, o)
}

func (pg *Pg) checkInsertNewOrderResults(res models.OrderInsertResult, o models.OrderNew) error {
	log.Println("pg.InsertOrder results:", o, res)
	if isForeignKeyViolation(res.Err) {
		log.Println("pg.InsertOrder error: incorrect user_id FK", o, res)
		return ErrUserUnknown
	}
	if res.Err != nil {
		log.Printf("\n pg.InsertOrder unexpected error: %#v \n\n", res.Err)
		return ErrUnxpectedError
	}
	if res.User != o.User {
		log.Println("другой юзер уже загрузил этот номер заказа", o, res)
		return ErrOrderWasPostedByAnotherUser
	}
	if res.WasConflict {
		log.Println("номер заказа уже был загружен этим юзером", o, res)
		return ErrOrderWasPostedByThisUser
	}
	log.Println("успешно принят новый заказ", o, res)
	return nil
}
