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

const (
	connMaxLifetime = 5 * time.Minute
	maxOpenConns    = 10
	maxIdleConns    = 5
)

func Connect(ctx context.Context, cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", buildDSN(cfg))
	if err != nil {
		return nil, err
	}

	configurePool(db)

	if err := db.PingContext(ctx); err != nil {
		closeErr := db.Close()
		return nil, errors.Join(err, closeErr)
	}

	return db, nil
}

func buildDSN(cfg config.DatabaseConfig) string {
	return (&mysql.Config{
		User:                 cfg.User,
		Passwd:               cfg.Password,
		Net:                  "tcp",
		Addr:                 net.JoinHostPort(cfg.Host, cfg.Port),
		DBName:               cfg.Name,
		ParseTime:            true,
		Collation:            "utf8mb4_unicode_ci",
		AllowNativePasswords: true,
		Params: map[string]string{
			"charset": "utf8mb4",
		},
	}).FormatDSN()
}

func configurePool(db *sql.DB) {
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
}
