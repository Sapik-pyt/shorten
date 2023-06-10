package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pressly/goose/v3"
)

// Создания pool с postgres
func ConnectToDb(ctx context.Context) (*pgxpool.Pool, error) {
	connStr := newConnectionString()

	pool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}

	err = migrate(connStr)
	if err != nil {
		return nil, fmt.Errorf("migrations failed to apply: %s", err.Error())
	}

	return pool, nil
}

// Строка для конекта к postgres
func newConnectionString() string {
	str := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%s",
		"postgres",
		"postgres",
		"123",
		"localhost",
		"5432",
		"postgres",
		"5")

	return str
}

// Обновление миграций 
func migrate(connStr string) error {
	db, err := goose.OpenDBWithDriver("pgx", connStr)
	if err != nil {
		return fmt.Errorf("goose: failed to open DB: %v\n", err)
	}
	defer db.Close()

	if err := goose.Up(db, "./migrations"); err != nil {
		return err
	}

	return nil
}
