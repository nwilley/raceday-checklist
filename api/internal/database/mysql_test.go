package database

import (
	"database/sql"
	"strings"
	"testing"

	"raceday-checklist/api/internal/config"

	"github.com/go-sql-driver/mysql"
)

func TestBuildDSN(t *testing.T) {
	dsn := buildDSN(config.DatabaseConfig{
		Host:     "db.example.test",
		Port:     "3307",
		User:     "raceday",
		Password: "secret",
		Name:     "checklists",
	})

	parsed, err := mysql.ParseDSN(dsn)
	if err != nil {
		t.Fatalf("parse dsn: %v", err)
	}

	if parsed.User != "raceday" {
		t.Fatalf("expected user raceday, got %q", parsed.User)
	}
	if parsed.Passwd != "secret" {
		t.Fatalf("expected password secret, got %q", parsed.Passwd)
	}
	if parsed.Net != "tcp" {
		t.Fatalf("expected tcp network, got %q", parsed.Net)
	}
	if parsed.Addr != "db.example.test:3307" {
		t.Fatalf("expected address db.example.test:3307, got %q", parsed.Addr)
	}
	if parsed.DBName != "checklists" {
		t.Fatalf("expected database checklists, got %q", parsed.DBName)
	}
	if !parsed.ParseTime {
		t.Fatal("expected parseTime to be enabled")
	}
	if parsed.Collation != "utf8mb4_unicode_ci" {
		t.Fatalf("expected collation utf8mb4_unicode_ci, got %q", parsed.Collation)
	}
	if !strings.Contains(dsn, "charset=utf8mb4") {
		t.Fatalf("expected dsn to include charset=utf8mb4, got %q", dsn)
	}
}

func TestConfigurePool(t *testing.T) {
	db, err := sql.Open("mysql", buildDSN(config.DatabaseConfig{
		Host: "127.0.0.1",
		Port: "3306",
		User: "raceday",
		Name: "checklists",
	}))
	if err != nil {
		t.Fatalf("open db handle: %v", err)
	}
	defer db.Close()

	configurePool(db)

	stats := db.Stats()
	if stats.MaxOpenConnections != maxOpenConns {
		t.Fatalf("expected max open connections %d, got %d", maxOpenConns, stats.MaxOpenConnections)
	}
}
