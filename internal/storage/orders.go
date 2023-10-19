package storage

import (
	"log"
	"time"

	"github.com/Kanbenn/gophermart/internal/models"
)

func (pg *Pg) InsertOrder(o models.OrderInsert) error {
	q := `
	INSERT INTO orders (number, user_id, uploaded_at) VALUES($1, $2, $3)
	ON CONFLICT (number) DO UPDATE SET updated_at = NOW()
	RETURNING user_id, (updated_at > created_at) AS was_conflict`

	o.UploadedAt = time.Now().Format(time.RFC3339)
	res := models.OrderInsertResult{}
	row := pg.Sqlx.QueryRowx(q, o.Number, o.UID, o.UploadedAt)
	res.Err = row.Scan(&res.UID, &res.WasConflict)
	return checkResults(res, o)
}

func checkResults(res models.OrderInsertResult, o models.OrderInsert) error {
	log.Println("pg.InsertOrder results:", o, res)
	if res.Err != nil {
		log.Printf("\n pg.InsertOrder unexpected error: %#v \n\n", res.Err)
		return ErrUnxpectedError
	}
	if res.UID != o.UID {
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
