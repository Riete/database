package conn

import (
	"context"
	"database/sql"
	"time"
)

type Conn struct {
	db *sql.DB
}

func (c Conn) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return c.db.QueryContext(ctx, query, args...)
}

func (c Conn) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}

func (c Conn) ExecWithTransaction(ctx context.Context, query string, args ...any) (sql.Result, error) {
	var t *sql.Tx
	var r sql.Result
	var err error
	t, err = c.db.BeginTx(ctx, nil)
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

func (c Conn) QueryWithTransaction(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	var t *sql.Tx
	var rows *sql.Rows
	var err error
	t, err = c.db.BeginTx(ctx, nil)
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

func (c Conn) Transaction(ctx context.Context, f func(*sql.Tx) error) error {
	t, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err = f(t); err != nil {
		return t.Rollback()
	}
	return t.Commit()
}

func (c Conn) DB() *sql.DB {
	return c.db
}

func (c *Conn) Close() error {
	return c.db.Close()
}

func SetDBConn(db *sql.DB, maxOpenConns, maxIdleConns int, maxConnLifetime time.Duration) {
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

func New(db *sql.DB) *Conn {
	return &Conn{db: db}
}
