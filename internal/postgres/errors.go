package postgres

import (
	"log"

	"github.com/Kanbenn/gophermart/internal/models"
	"github.com/lib/pq"
)

type orderInsertResult struct {
	User        int
	WasConflict bool
	SameUser    bool
	Number      string
	Err         error
}

func isNotUniqueInsert(e error) bool {
	if err, ok := e.(*pq.Error); ok && err.Code == "23505" {
		return true
	}
	return false
}

func isNotEnoughMinerals(e error) bool {
	if err, ok := e.(*pq.Error); ok && err.Code == "23514" {
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

func (pg *Pg) checkInsertOrderResults(res orderInsertResult) error {
	if !res.SameUser {
		log.Println("другой юзер уже загрузил этот номер заказа", res)
		return models.ErrOrderWasPostedByAnotherUser
	}
	if res.WasConflict && res.SameUser {
		log.Println("номер заказа уже был загружен этим юзером", res)
		return models.ErrOrderWasPostedByThisUser
	}
	if isForeignKeyViolation(res.Err) {
		log.Println("pg.InsertOrder error: incorrect user_id FK")
		return models.ErrUserUnknown
	}
	if res.Err != nil && !isForeignKeyViolation(res.Err) {
		log.Printf("\n pg.InsertOrder unexpected error: %#v \n\n", res.Err)
		return models.ErrUnxpectedError
	}
	log.Println("успешно принят новый заказ", res.Number)
	return nil
}

func (pg *Pg) checkInserOrderWithBonusErr(err error) error {
	if isNotUniqueInsert(err) {
		log.Println("заказ c таким номером уже существует")
		return models.ErrOrderWasPostedByThisUser
	}
	log.Println("неожиданная ошибка при добавлении нового заказа withdrawal:")
	log.Printf("%#v \n\n", err)
	return models.ErrUnxpectedError
}

func (pg *Pg) checkUpdateUserBalanceErr(err error) error {
	if isNotEnoughMinerals(err) {
		log.Println("ошибка при добавлении нового заказа withdrawals")
		log.Println("недостаточно бонусных баллов на балансе пользователя:")
		return models.ErrNotEnoughMinerals
	}
	log.Println("неожиданная ошибка при списании баланса пользователя для заказа withdrawals:")
	log.Printf("%#v \n\n", err)
	return models.ErrUnxpectedError
}
