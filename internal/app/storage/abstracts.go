package storage

import "errors"

var (
	ErrURLMappingNotFound = errors.New("URL not found")
)

type AbstractStorage interface {
	GetURLByShortLink(shortLink string) (string, error)
	AddShortURL(shortLink string, fullLink string) error
}