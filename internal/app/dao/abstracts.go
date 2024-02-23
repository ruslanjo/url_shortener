package dao

import "errors"

var (
	ErrURLMappingNotFound = errors.New("URL not found")
)

type AbstractDAO interface {
	GetURLByShortLink(shortLink string) (string, error)
	AddShortURL(shortLink string, fullLink string) error
}