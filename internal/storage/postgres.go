package storage

import (
	"log"

	"github.com/jmoiron/sqlx" // needs a pg driver like github.com/lib/pq
	"github.com/lib/pq"

	"github.com/Kanbenn/gophermart/internal/config"
	"github.com/Kanbenn/gophermart/internal/models"
)

const (
	createTables = `
	DROP TABLE IF EXISTS gm_users;
	CREATE TABLE gm_users (
		id SERIAL PRIMARY KEY,
		login TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		updated_at TIMESTAMP without time zone DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS ulogin ON gm_users(login);`
)

type Pg struct {
	Sqlx *sqlx.DB
	Cfg  config.Config
}

func NewPostgres(cfg config.Config) *Pg {
	x, err := sqlx.Open("postgres", cfg.PgConnStr)
	if err != nil {
		log.Fatal("error at connecting to Postgres:", cfg.PgConnStr, err)
	}
	if err := x.Ping(); err != nil {
		log.Fatal("error at pinging the Postgres db:", cfg.PgConnStr, err)
	}
	if _, err := x.Exec(createTables); err != nil {
		log.Fatal("error at creating db-tables:", cfg.PgConnStr, err)
	}
	return &Pg{x, cfg}
}

func (pg *Pg) Close() error {
	return pg.Sqlx.Close()
}

func (pg *Pg) InsertUser(user models.User) (uid int, err error) {
	q := `INSERT INTO gm_users(login, password) VALUES($1, $2)
		  RETURNING id`
	err = pg.Sqlx.QueryRowx(q, user.Login, user.Password).Scan(&uid)
	if isNotUnique(err) {
		return uid, ErrLoginNotUnique
	}
	return uid, err
}

func (pg *Pg) GetUser(login string) (user models.User) {
	q := `SELECT id, login, password FROM gm_users WHERE login = $1`
	_ = pg.Sqlx.Get(&user, q, login)
	return user
}

func isNotUnique(e error) bool {
	if err, ok := e.(*pq.Error); ok && err.Code == "23505" {
		return true
	}
	return false
}
