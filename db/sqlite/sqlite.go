package sqlite

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Option func(map[string]string)

func WithReadOnly() Option {
	return func(m map[string]string) {
		m["mode"] = "ro"
	}
}

func WithQueryOnly() Option {
	return func(m map[string]string) {
		m["_pragma.query_only"] = "query_only(1)"
	}
}

func WithForeignKeys() Option {
	return func(m map[string]string) {
		m["_pragma.foreign_keys"] = "foreign_keys(1)"
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
		m["_pragma.journal_mode"] = "journal_mode(" + strings.ToUpper(mode) + ")"
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
	query := make(url.Values)
	for key, value := range params {
		if strings.HasPrefix(key, "_pragma.") {
			key = "_pragma"
		}
		query.Add(key, value)
	}
	return fmt.Sprintf("%s?%s", dsn, query.Encode())
}

func NewDB(path string, options ...Option) (*sql.DB, error) {
	return sql.Open("sqlite", FormatDSN(path, options...))
}

func NewGormDB(path string, options ...Option) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(FormatDSN(path, options...)))
}
