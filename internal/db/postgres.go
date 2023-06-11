package db

import (
	"context"
	"fmt"

	_ "github.com/Sapik-pyt/shorten/internal/db/migrations"
	"github.com/Sapik-pyt/shorten/internal/logging"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// Создание connection pool
func ConnectToDb(ctx context.Context) (*pgxpool.Pool, error) {
	logging.Logger.Info("connecting to db")

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
		"db",
		"5432",
		"postgres",
		"5")

	return str
}

// Мигрирование БД
func migrate(connStr string) error {
	logging.Logger.Info("migration db")

	db, err := goose.OpenDBWithDriver("postgres", connStr)
	if err != nil {
		return fmt.Errorf("goose: failed to open DB: %s", err.Error())
	}
	defer db.Close()

	if err := goose.Up(db, "./internal/db/migrations"); err != nil {
		return err
	}

	return nil
}
