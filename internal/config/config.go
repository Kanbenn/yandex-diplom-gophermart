package config

import (
	"flag"
	"os"
	"time"
)

type Config struct {
	Addr                  string
	PgConnStr             string
	AccrualLink           string
	ProcessedAtTimeFormat string
	NewBonusOrderStatus   string
	FinishedOrderStatuses []string
}

func NewFromFlagsAndEnvs() Config {
	c := Config{}
	flag.StringVar(&c.Addr, "a", "localhost:8080", "address:port for Gophermart to listen on")
	flag.StringVar(&c.PgConnStr, "d", "", "Postgres connection string for storage")
	flag.StringVar(&c.AccrualLink, "r", "http://localhost:9090", "address:port for Accrual to listen on")
	flag.Parse()

	if val := os.Getenv("RUN_ADDRESS"); val != "" {
		c.Addr = val
	}
	if val := os.Getenv("DATABASE_URI"); val != "" {
		c.PgConnStr = val
	}
	if val := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); val != "" {
		c.AccrualLink = val
	}

	c.ProcessedAtTimeFormat = time.RFC3339
	c.NewBonusOrderStatus = "PROCESSED"
	c.FinishedOrderStatuses = []string{"PROCESSED", "INVALID"}

	return c
}
