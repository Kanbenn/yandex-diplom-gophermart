package storage

import (
	"errors"
	"log"

	"github.com/Kanbenn/gophermart/internal/models"
	"github.com/lib/pq"
)

var (
	// user errors
	ErrLoginNotUnique = errors.New("this login is already taken")
	ErrUserUnknown    = errors.New("unknown user")
	ErrNotAuthorized  = errors.New("unauthorized")

	// insertOrder errors
	ErrOrderWasPostedByThisUser    = errors.New("this order number was already posted")
	ErrOrderWasPostedByAnotherUser = errors.New("this order number was already posted by other user")

	ErrUnxpectedError = errors.New("unexpected server error")
	ErrBadRequest     = errors.New("bad request")
)

func isNotUniqueInsert(e error) bool {
	if err, ok := e.(*pq.Error); ok && err.Code == "23505" {
		return true
	}
	return false
}

func isForeignKeyViolation(e error) bool {
	if err, ok := e.(*pq.Error); ok && err.Code == "23503" {
		return true
	}
	return false
}

func checkInsertOrderResults(res models.OrderInsertResult, o models.OrderInsert) error {
	log.Println("pg.InsertOrder results:", o, res)
	if isForeignKeyViolation(res.Err) { // TODO change sql to: user_id INTEGER REFERENCES users(id),
		log.Println("pg.InsertOrder error: incorrect user_id FK", o, res)
		return ErrUserUnknown
	}
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
