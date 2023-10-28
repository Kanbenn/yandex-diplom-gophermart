package postgres

import (
	"log"
	"time"

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

func New(cfg config.Config) *Pg {
	conn, err := sqlx.Open("postgres", cfg.PgConnStr)
	if err != nil {
		log.Fatal("error at connecting to Postgres:", cfg.PgConnStr, err)
	}

	pg := Pg{conn, cfg}

	if _, err := pg.Sqlx.Exec(createTables); err != nil {
		log.Fatal("error at creating db-tables:", pg.Cfg.PgConnStr, err)
	}
	return &pg
}

func (pg *Pg) Close() error {
	return pg.Sqlx.Close()
}

func (pg *Pg) statusForNewOrderWithBonus() string {
	// чтобы новые заказы withdrawal не попадали на запросы к Accrual
	return pg.Cfg.NewBonusOrderStatus
}

func (pg *Pg) timeForNewOrders() string {
	// по Тех Заданию, для генерации json:"processed_at"
	return time.Now().Format(pg.Cfg.ProcessedAtTimeFormat)
}
