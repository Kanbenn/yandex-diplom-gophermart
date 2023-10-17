package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr        string
	PgConnStr   string
	AccrualAddr string
}

func NewFromFlagsAndEnvs() Config {
	c := Config{}
	flag.StringVar(&c.Addr, "a", "localhost:8080", "address:port for Gophermart to listen on")
	flag.StringVar(&c.PgConnStr, "d", "", "Postgres connection string for storage")
	flag.StringVar(&c.AccrualAddr, "r", "localhost:9090", "address:port for Accrual to listen on")
	flag.Parse()

	if val := os.Getenv("RUN_ADDRESS"); val != "" {
		c.Addr = val
	}
	if val := os.Getenv("DATABASE_URI"); val != "" {
		c.PgConnStr = val
	}
	if val := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); val != "" {
		c.AccrualAddr = val
	}
	return c
}
