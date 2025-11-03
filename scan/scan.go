package scan

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type Scanner[T any] struct {
	rows    *sql.Rows
	db      *gorm.DB
	scanErr error
}

func (s *Scanner[T]) Columns() ([]string, error) {
	return s.rows.Columns()
}

func (s *Scanner[T]) ScanFirst() (*T, error) {
	if !s.rows.Next() {
		return nil, sql.ErrNoRows
	}
	scanTo := new(T)
	s.scanErr = s.db.ScanRows(s.rows, scanTo)
	return scanTo, s.scanErr
}

func (s *Scanner[T]) Scan(ctx context.Context) <-chan *T {
	ch := make(chan *T)
	go func() {
		defer close(ch)
		for s.rows.Next() {
			select {
			case <-ctx.Done():
				return
			default:
				scanTo := new(T)
				if s.scanErr = s.db.ScanRows(s.rows, scanTo); s.scanErr != nil {
					return
				}
				ch <- scanTo
			}
		}
	}()
	return ch
}

func (s *Scanner[T]) ScanToMap(ctx context.Context) <-chan map[string]string {
	columns, err := s.Columns()
	if err != nil {
		return nil
	}
	ch := make(chan map[string]string)
	go func() {
		defer close(ch)
		values := make([]sql.RawBytes, len(columns))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		for s.rows.Next() {
			select {
			case <-ctx.Done():
				return
			default:
				if s.scanErr = s.rows.Scan(scanArgs...); s.scanErr != nil {
					return
				}
				row := make(map[string]string, len(columns))
				for i, v := range values {
					if v == nil {
						row[columns[i]] = "null"
					} else {
						row[columns[i]] = string(v)
					}
				}
				ch <- row
			}
		}
	}()
	return ch
}

func (s *Scanner[T]) HasNextRS() bool {
	return s.rows.NextResultSet()
}

func (s *Scanner[T]) Error() error {
	if s.scanErr != nil {
		return s.scanErr
	}
	return s.rows.Err()
}

func (s *Scanner[T]) Close() error {
	return s.rows.Close()
}

func NewScanner[T any](rows *sql.Rows, db *gorm.DB) *Scanner[T] {
	return &Scanner[T]{rows: rows, db: db}
}
