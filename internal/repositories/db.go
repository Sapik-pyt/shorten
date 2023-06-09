package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Стркутура хранилища на основе БД
type dbStorage struct {
	pool *pgxpool.Pool
}

// Функция создания хранилища на основе БД
func NewDbStorage(pool *pgxpool.Pool) *dbStorage {
	return &dbStorage{
		pool: pool,
	}
}

// Метод для получения оригинальной ссылки из БД
func (db *dbStorage) Get(ctx context.Context, shortLink string) (*string, error) {

	query := `SELECT original_link FROM links WHERE short_link=$1`
	links := linkData{}

	res := db.pool.QueryRow(ctx, query, shortLink)
	err := res.Scan(&links.OriginalLink)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scanning result: %s", err.Error())
	}
	return &links.OriginalLink, nil
}

// Метод для сохранения данных в БД
func (db *dbStorage) Save(ctx context.Context, shortLink, originaLink string) error {
	query := `INSERT INTO links(short_link, original_link, created_at) VALUES($1, $2, $3)`
	links := linkData{
		ShortLink:    shortLink,
		OriginalLink: originaLink,
		CreatedAt:    time.Now(),
	}
	_, err := db.pool.Exec(ctx, query, links.ShortLink, links.OriginalLink, links.CreatedAt)
	if err != nil {
		return fmt.Errorf("inserting into db: %s", err.Error())
	}
	return nil
}

// Метод для проверки существования ссылки в БД
func (db *dbStorage) CheckExistance(ctx context.Context, shortLink string) (bool, error) {
	query := `SELECT original_link FROM links WHERE short_link=$1`
	links := linkData{}

	res := db.pool.QueryRow(ctx, query, shortLink)
	err := res.Scan(&links.OriginalLink)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("scanning result: %s", err.Error())
	}
	return true, nil

}

// Структура для считывания данных из БД
type linkData struct {
	ID           int
	ShortLink    string
	OriginalLink string
	CreatedAt    time.Time
}
