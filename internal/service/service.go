package service

import (
	"context"

	gen "github.com/Sapik-pyt/shorten/proto/gen"
)

type Repository interface {
	Get(ctx context.Context, shortLink string) (*string, error)
	Save(ctx context.Context, shortLink, originaLink string) error
	CheckExistance(ctx context.Context, shortLink string) (bool, error)
}

// Структура сервиса
type ShortenService struct {
	gen.UnimplementedShortenServer
	repository Repository
}

// Функция создания сервиса
func NewShortenService(repository Repository) *ShortenService {
	return &ShortenService{
		repository: repository,
	}
}
