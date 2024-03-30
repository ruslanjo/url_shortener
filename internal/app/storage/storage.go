package storage

import (
	"context"
	"errors"

	"github.com/ruslanjo/url_shortener/internal/app/storage/models"
)

var (
	ErrURLMappingNotFound = errors.New("URL not found")
	ErrIntegityViolation  = errors.New("db constrains violation")
	ErrEntityDeleted      = errors.New("entity was deleted")
)

type Storage interface {
	GetURLByShortLink(shortLink string) (string, error)
	AddShortURL(shortLink string, fullLink string, UUID string) error
	SaveURLBatched(ctx context.Context, data []models.URLBatch, UUID string) error
	PingContext(ctx context.Context) error
	GetUserURLs(UUID string) ([]models.URL, error)
	DeleteURLs(ctx context.Context, shortURLs []string, userID string) error
}
