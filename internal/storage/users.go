package storage

import (
	"github.com/Kanbenn/gophermart/internal/models"
)

func (pg *Pg) InsertUser(user models.UserInsert) (uid int, err error) {
	q := `INSERT INTO users(login, password) VALUES($1, $2)
		  RETURNING id`
	err = pg.Sqlx.QueryRowx(q, user.Login, user.Password).Scan(&uid)
	if isNotUniqueInsert(err) {
		return uid, ErrLoginNotUnique
	}
	return uid, err
}

func (pg *Pg) GetUser(login string) (user models.UserInsert) {
	q := `SELECT id, login, password FROM users WHERE login = $1`
	_ = pg.Sqlx.Get(&user, q, login)
	return user
}
