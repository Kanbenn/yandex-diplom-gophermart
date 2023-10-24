package storage

import (
	"log"
	"time"

	"github.com/Kanbenn/gophermart/internal/models"
)

func (pg *Pg) InsertOrder(o models.Order) error {
	q := `
	INSERT INTO orders (number, user_id, time) VALUES($1, $2, $3)
	ON CONFLICT (number) DO UPDATE SET updated_at = NOW()
	RETURNING user_id, (updated_at > created_at) AS was_conflict`

	o.Time = timeForNewOrders()
	result := models.OrderInsertResult{}
	result.Err = pg.Sqlx.QueryRowx(q, o.Number, o.User, o.Time).
		Scan(&result.User, &result.WasConflict)
	return pg.checkInsertOrderResults(result, o)
}

func (pg *Pg) InsertOrderWithBonus(o models.Order) error {
	tx, err := pg.Sqlx.Beginx()
	if err != nil {
		log.Printf("\n pg.InsertOrderWithBonus unexpected tx error: %#v \n\n", err)
		return ErrUnxpectedError
	}
	defer tx.Rollback()

	qo := `
	INSERT INTO orders (number, user_id, status, sum, time) 
	VALUES(:number, :user, :status, :sum, :time)`
	o.Time = timeForNewOrders()
	o.Status = statusForNewOrderWithBonus()
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

func (pg *Pg) checkInsertOrderResults(res models.OrderInsertResult, o models.Order) error {
	log.Println("pg.InsertOrder results:", o, res)
	if isForeignKeyViolation(res.Err) {
		log.Println("pg.InsertOrder error: incorrect user_id FK", o, res)
		return ErrUserUnknown
	}
	if res.User != o.User {
		log.Println("другой юзер уже загрузил этот номер заказа", o, res)
		return ErrOrderWasPostedByAnotherUser
	}
	if res.WasConflict {
		log.Println("номер заказа уже был загружен этим юзером", o, res)
		return ErrOrderWasPostedByThisUser
	}
	if res.Err != nil {
		log.Printf("\n pg.InsertOrder unexpected error: %#v \n\n", res.Err)
		return ErrUnxpectedError
	}
	log.Println("успешно принят новый заказ", o)
	return nil
}

func statusForNewOrderWithBonus() string {
	// чтобы такие новые заказы (withdrawal) не попадали на запросы к Accrual
	return "PROCESSED"
}

func timeForNewOrders() string {
	// по Тех Заданию, для генерации json:"processed_at"
	return time.Now().Format(time.RFC3339)
}

func checkUpdateUserBalanceErr(err error) error {
	if isNotEnoughMinerals(err) {
		log.Println("ошибка при добавлении нового заказа withdrawals")
		log.Println("недостаточно бонусных баллов на балансе пользователя:")
		return ErrNotEnoughMinerals
	}
	log.Println("неожиданная ошибка при списании баланса пользователя для заказа withdrawals:")
	log.Printf("%#v \n\n", err)
	return ErrUnxpectedError
}

func checkInserOrderWithBonusErr(err error) error {
	if isNotUniqueInsert(err) {
		log.Println("заказ c таким номером уже существует")
		return ErrOrderWasPostedByThisUser
	}
	log.Println("неожиданная ошибка при добавлении нового заказа withdrawal:")
	log.Printf("%#v \n\n", err)
	return ErrUnxpectedError
}
