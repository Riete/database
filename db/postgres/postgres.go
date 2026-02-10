package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSlMode  bool
}

func NewConn(host, port, username, password, dbName string) (*pgx.Conn, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbName)
	return pgx.Connect(context.Background(), dsn)
}

func NewDB(host, port, username, password, dbName string) (*sql.DB, *gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		host,
		username,
		password,
		dbName,
		port,
	)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}))
	if err == nil {
		var sqlDB *sql.DB
		if sqlDB, err = gormDB.DB(); err == nil {
			return sqlDB, gormDB, nil
		}
	}
	return nil, nil, err
}
