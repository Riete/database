package mysql

import (
	"context"
	"database/sql"
	"net"
	"sync"
	"time"

	"golang.org/x/net/proxy"
	"gorm.io/gorm"

	"github.com/go-sql-driver/mysql"
	gmysql "gorm.io/driver/mysql"
)

var registered = make(map[string]bool)
var mu sync.Mutex

func WithSocks5Proxy(addr string, auth *proxy.Auth) mysql.Option {
	return func(config *mysql.Config) error {
		mu.Lock()
		defer mu.Unlock()
		proxyNet := "socks5://" + addr
		if !registered[proxyNet] {
			p, _ := proxy.SOCKS5("tcp", addr, auth, proxy.Direct)
			mysql.RegisterDialContext(proxyNet, func(ctx context.Context, addr string) (net.Conn, error) {
				return p.(proxy.ContextDialer).DialContext(ctx, "tcp", addr)
			})
			registered[proxyNet] = true
		}
		config.Net = proxyNet
		return nil
	}
}

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
	config.Net = "tcp"
	return config
}

func NewConfig(options ...mysql.Option) *mysql.Config {
	config := NewDefaultConfig()
	_ = config.Apply(options...)
	return config
}

func NewDB(config *mysql.Config) (*sql.DB, error) {
	return sql.Open("mysql", config.FormatDSN())
}

func NewGormDB(db *sql.DB) (*gorm.DB, error) {
	return gorm.Open(gmysql.New(gmysql.Config{Conn: db}))
}
