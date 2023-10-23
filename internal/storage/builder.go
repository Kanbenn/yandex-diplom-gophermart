package storage

import (
	"log"

	"github.com/jmoiron/sqlx" // needs a pg driver like github.com/lib/pq

	"github.com/Kanbenn/gophermart/internal/config"
)

const (
	createTables = `
	DROP TABLE IF EXISTS users, orders;

	CREATE TABLE IF NOT EXISTS users (
		id 		   SERIAL PRIMARY KEY,
		login      TEXT NOT NULL UNIQUE,
		password   TEXT NOT NULL,
		updated_at TIMESTAMP without time zone DEFAULT NOW(),
		balance    DECIMAL(12,2) DEFAULT 0,
		withdrawn  DECIMAL(12,2) DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS ulogin ON users(login);
	
	CREATE TABLE IF NOT EXISTS orders (
		id 		    SERIAL PRIMARY KEY,
		number      TEXT NOT NULL UNIQUE,
		user_id     INTEGER REFERENCES users(id),
		status      TEXT DEFAULT 'NEW',
		bonus		DECIMAL(12,2) DEFAULT 0,
		sum  		DECIMAL(12,2) DEFAULT 0,
		time		TEXT,
		created_at  TIMESTAMP without time zone DEFAULT NOW(),
		updated_at  TIMESTAMP without time zone DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS onumber ON orders(number);`
)

type Pg struct {
	Sqlx *sqlx.DB
	Cfg  config.Config
}

func NewPostgres(cfg config.Config) *Pg {
	conn, err := sqlx.Open("postgres", cfg.PgConnStr)
	if err != nil {
		log.Fatal("error at connecting to Postgres:", cfg.PgConnStr, err)
	}
	if err := conn.Ping(); err != nil {
		log.Fatal("error at pinging the Postgres db:", cfg.PgConnStr, err)
	}
	if _, err := conn.Exec(createTables); err != nil {
		log.Fatal("error at creating db-tables:", cfg.PgConnStr, err)
	}

	return &Pg{conn, cfg}
}

func (pg *Pg) Close() error {
	return pg.Sqlx.Close()
}
