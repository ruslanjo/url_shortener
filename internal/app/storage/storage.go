package storage

import (
	"context"
	"errors"

	"github.com/ruslanjo/url_shortener/internal/app/storage/models"
)

var (
	ErrURLMappingNotFound = errors.New("URL not found")
	ErrIntegityViolation = errors.New("db constrains violation")
)

type Storage interface {
	GetURLByShortLink(shortLink string) (string, error)
	AddShortURL(shortLink string, fullLink string) error
	SaveURLBatched(ctx context.Context, data []models.URLBatch) error
	PingContext(ctx context.Context) error
}
