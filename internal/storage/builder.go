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
		balance    DECIMAL(12,2) CHECK (balance >= 0) DEFAULT 0,
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
	CREATE INDEX IF NOT EXISTS onumber ON orders(number);
	CREATE INDEX IF NOT EXISTS ostatus ON orders(status);`
)

type Pg struct {
	Sqlx *sqlx.DB
	Cfg  config.Config
}

func NewInPostgres(cfg config.Config, conn *sqlx.DB) *Pg {
	return &Pg{conn, cfg}
}

func (pg *Pg) CreateTables() {
	if _, err := pg.Sqlx.Exec(createTables); err != nil {
		log.Fatal("error at creating db-tables:", pg.Cfg.PgConnStr, err)
	}
}

func (pg *Pg) Close() error {
	return pg.Sqlx.Close()
}
