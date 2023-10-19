package storage

import (
	"github.com/lib/pq"

	"github.com/Kanbenn/gophermart/internal/models"
)

func (pg *Pg) InsertUser(user models.UserInsert) (uid int, err error) {
	q := `INSERT INTO users(login, password) VALUES($1, $2)
		  RETURNING id`
	err = pg.Sqlx.QueryRowx(q, user.Login, user.Password).Scan(&uid)
	if isNotUnique(err) {
		return uid, ErrLoginNotUnique
	}
	return uid, err
}

func (pg *Pg) GetUser(login string) (user models.UserInsert) {
	q := `SELECT id, login, password FROM users WHERE login = $1`
	_ = pg.Sqlx.Get(&user, q, login)
	return user
}

func isNotUnique(e error) bool {
	if err, ok := e.(*pq.Error); ok && err.Code == "23505" {
		return true
	}
	return false
}
