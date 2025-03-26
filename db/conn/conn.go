package conn

import (
	"database/sql"
	"time"
)

func SetConn(db *sql.DB, maxOpenConns, maxIdleConns int, maxConnLifetime time.Duration) {
	if maxOpenConns <= 0 {
		maxOpenConns = 20
	}
	if maxIdleConns <= 0 {
		maxIdleConns = 20
	}
	if maxConnLifetime <= 0 {
		maxConnLifetime = 3 * time.Minute
	}
	db.SetConnMaxLifetime(maxConnLifetime)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
}
