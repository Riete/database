package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Option func(*MySQL)

func WithMaxConn(max int) Option {
	return func(m *MySQL) {
		m.maxConn = max
	}
}

func WithMaxConnLifetime(t time.Duration) Option {
	return func(m *MySQL) {
		m.maxConnLifetime = t
	}
}

var DefaultConfig = &mysql.Config{
	Timeout:         10 * time.Second,
	Loc:             time.Local,
	Params:          map[string]string{"charset": "utf8mb4"},
	ParseTime:       true,
	ClientFoundRows: true,
}

type MySQL struct {
	db              *sql.DB
	maxConn         int
	maxConnLifetime time.Duration
	config          *mysql.Config
}

func (m *MySQL) open() error {
	var err error
	m.db, err = sql.Open("mysql", m.config.FormatDSN())
	if err != nil {
		return err
	}
	m.db.SetConnMaxLifetime(m.maxConnLifetime)
	m.db.SetMaxOpenConns(m.maxConn)
	m.db.SetMaxIdleConns(m.maxConn)
	return nil
}

func (m MySQL) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, query, args...)
}

func (m MySQL) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return m.db.ExecContext(ctx, query, args...)
}

func (m MySQL) ExecWithTransaction(ctx context.Context, query string, args ...any) (sql.Result, error) {
	var t *sql.Tx
	var r sql.Result
	var err error
	t, err = m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	if r, err = t.ExecContext(ctx, query, args...); err != nil {
		if rbErr := t.Rollback(); rbErr != nil {
			return nil, rbErr
		}
		return nil, err
	}
	return r, t.Commit()
}

func (m MySQL) QueryWithTransaction(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	var t *sql.Tx
	var rows *sql.Rows
	var err error
	t, err = m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	if rows, err = t.QueryContext(ctx, query, args...); err != nil {
		if rbErr := t.Rollback(); rbErr != nil {
			return nil, rbErr
		}
		return nil, err
	}
	return rows, t.Commit()

}

func (m MySQL) Transaction(ctx context.Context, f func(*sql.Tx) error) error {
	t, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err = f(t); err != nil {
		return t.Rollback()
	}
	return t.Commit()
}

func (m MySQL) DSN() string {
	return m.config.FormatDSN()
}

func (m MySQL) SqlDB() *sql.DB {
	return m.db
}

func (m MySQL) GormDB() (*gorm.DB, error) {
	return gorm.Open(gmysql.New(gmysql.Config{Conn: m.db}))
}

func (m *MySQL) Close() error {
	return m.db.Close()
}

func New(config *mysql.Config, options ...Option) (*MySQL, error) {
	m := &MySQL{
		maxConn:         20,
		maxConnLifetime: 3 * time.Minute,
		config:          config,
	}
	for _, option := range options {
		option(m)
	}
	return m, m.open()
}
