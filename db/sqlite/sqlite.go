package sqlite

import (
	"database/sql"
	"fmt"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/mattn/go-sqlite3"
)

type Option func(map[string]string)

func WithReadOnly() Option {
	return func(m map[string]string) {
		m["mode"] = "ro"
	}
}

func WithQueryOnly() Option {
	return func(m map[string]string) {
		m["_query_only"] = "true"
	}
}

func WithInMemory() Option {
	return func(m map[string]string) {
		m["mode"] = "memory"
	}
}

func WithCacheShared() Option {
	return func(m map[string]string) {
		m["cache"] = "shared"
	}
}

// WithJournalMode DELETE | TRUNCATE | PERSIST | MEMORY | WAL | OFF
func WithJournalMode(mode string) Option {
	return func(m map[string]string) {
		m["_journal"] = strings.ToUpper(mode)
	}
}

func FormatDSN(name string, options ...Option) string {
	dsn := fmt.Sprintf("file:%s", name)
	if len(options) == 0 {
		return dsn
	}
	params := make(map[string]string)
	for _, option := range options {
		option(params)
	}
	var param []string
	for k, v := range params {
		if v != "" {
			param = append(param, fmt.Sprintf("%s=%s", k, v))
		} else {
			param = append(param, k)
		}
	}
	return fmt.Sprintf("%s?%s", dsn, strings.Join(param, "&"))
}

func NewDB(path string, options ...Option) (*sql.DB, error) {
	return sql.Open("sqlite3", FormatDSN(path, options...))
}

func NewGormDB(db *sql.DB) (*gorm.DB, error) {
	return gorm.Open(sqlite.New(sqlite.Config{Conn: db}))
}
