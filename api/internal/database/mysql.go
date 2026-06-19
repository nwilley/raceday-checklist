package database

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"time"

	"raceday-checklist/api/internal/config"

	"github.com/go-sql-driver/mysql"
)

func Connect(ctx context.Context, cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := (&mysql.Config{
		User:      cfg.User,
		Passwd:    cfg.Password,
		Net:       "tcp",
		Addr:      net.JoinHostPort(cfg.Host, cfg.Port),
		DBName:    cfg.Name,
		ParseTime: true,
		Params: map[string]string{
			"charset": "utf8mb4",
		},
	}).FormatDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err := db.PingContext(ctx); err != nil {
		closeErr := db.Close()
		return nil, errors.Join(err, closeErr)
	}

	return db, nil
}
