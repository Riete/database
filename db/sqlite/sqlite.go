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

// WithAuth go build -tags "sqlite_userauth" add such tags when compile to enable user auth
func WithAuth(username, password string) Option {
	return func(m map[string]string) {
		m["_auth"] = ""
		m["_auth_user"] = username
		m["_auth_pass"] = password
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

func formatDSN(name string, options ...Option) string {
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

func NewDB(name string, options ...Option) (*sql.DB, error) {
	return sql.Open("sqlite3", formatDSN(name, options...))
}

func NewGormDB(name string, options ...Option) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(formatDSN(name, options...)))
}
