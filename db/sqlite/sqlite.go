package sqlite

import (
	"database/sql"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB(path string) (*sql.DB, error) {
	return sql.Open("sqlite3", path)
}

func NewGormDB(path string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(path))
}
