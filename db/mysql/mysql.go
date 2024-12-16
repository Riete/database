package mysql

import (
	"database/sql"
	"time"

	"gorm.io/gorm"

	"github.com/go-sql-driver/mysql"
	gmysql "gorm.io/driver/mysql"
)

func WithAddr(host, port string) mysql.Option {
	return func(config *mysql.Config) error {
		config.Addr = host + ":" + port
		return nil
	}
}

func WithAuth(username, password string) mysql.Option {
	return func(config *mysql.Config) error {
		config.User = username
		config.Passwd = password
		return nil
	}
}

func WithDBName(dbname string) mysql.Option {
	return func(config *mysql.Config) error {
		config.DBName = dbname
		return nil
	}
}

func WithTimeout(timeout time.Duration) mysql.Option {
	return func(config *mysql.Config) error {
		config.Timeout = timeout
		return nil
	}
}

func NewDefaultConfig() *mysql.Config {
	config := mysql.NewConfig()
	config.Loc = time.Local
	config.Timeout = 10 * time.Second
	config.Params = map[string]string{"charset": "utf8mb4"}
	config.ParseTime = true
	return config
}

func NewConfig(options ...mysql.Option) (*mysql.Config, error) {
	config := NewDefaultConfig()
	return config, config.Apply(options...)
}

func NewDB(config *mysql.Config) (*sql.DB, error) {
	return sql.Open("mysql", config.FormatDSN())
}

func NewGormDB(db *sql.DB) (*gorm.DB, error) {
	return gorm.Open(gmysql.New(gmysql.Config{Conn: db}))
}
