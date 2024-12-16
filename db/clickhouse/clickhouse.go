package clickhouse

import (
	"database/sql"
	"net/url"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	gclickhouse "gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type Option func(options *clickhouse.Options)

func WithAddr(host, port string) Option {
	return func(options *clickhouse.Options) {
		options.Addr = []string{host + ":" + port}
	}
}

func WithAuth(username, password string) Option {
	return func(options *clickhouse.Options) {
		options.Auth.Username = username
		options.Auth.Password = password
	}
}

func WithDatabase(database string) Option {
	return func(options *clickhouse.Options) {
		options.Auth.Database = database
	}
}

func WithMaxOpenConns(n int) Option {
	return func(options *clickhouse.Options) {
		options.MaxOpenConns = n
	}
}

func WithMaxIdleConns(n int) Option {
	return func(options *clickhouse.Options) {
		options.MaxIdleConns = n
	}
}

func WithConnMaxLifetime(n time.Duration) Option {
	return func(options *clickhouse.Options) {
		options.ConnMaxLifetime = n
	}
}

func WithCompression(c *clickhouse.Compression) Option {
	return func(options *clickhouse.Options) {
		options.Compression = c
	}
}

func WithHTTPProxyURL(url *url.URL) Option {
	return func(options *clickhouse.Options) {
		options.HTTPProxyURL = url
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(options *clickhouse.Options) {
		options.DialTimeout = timeout
	}
}

func NewOptions(options ...Option) *clickhouse.Options {
	o := &clickhouse.Options{}
	for _, option := range options {
		option(o)
	}
	return o
}

func NewDB(option *clickhouse.Options) *sql.DB {
	return clickhouse.OpenDB(option)
}

func NewConn(option *clickhouse.Options) (clickhouse.Conn, error) {
	return clickhouse.Open(option)
}

func NewGormDB(db *sql.DB) (*gorm.DB, error) {
	return gorm.Open(gclickhouse.New(gclickhouse.Config{Conn: db}))
}
