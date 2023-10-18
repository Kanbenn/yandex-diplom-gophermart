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
	RETURNING user_id, (updated_at > created_at) AS conflict`

	o.UploadedAt = time.Now().Format(time.RFC3339)
	res := models.OrderInsertResult{}
	row := pg.Sqlx.QueryRowx(q, o.Number, o.UID, o.UploadedAt)
	res.Err = row.Scan(&res.UID, &res.IsConflict)
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
	if res.IsConflict {
		log.Println("номер заказа уже был загружен этим юзером", o, res)
		return ErrOrderWasPostedByThisUser
	}
	log.Println("успешно принят новый заказ", o, res)
	return nil
}

// err = pg.Sqlx.QueryRowx(q, user.Login, user.Password).Scan(&uid)
// if isNotUnique(err) {
// 	return uid, ErrLoginNotUnique
// }
// return uid, err
// _, err := pg.Sqlx.NamedExec(q, &user)
// if err != nil {
// 	if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
// 		log.Println("pg.InsertUser: login already taken", user)
// 		return ErrLoginTaken
// 	}
// 	log.Println("pg.InsertUser unknown error:", err, user)
// }
