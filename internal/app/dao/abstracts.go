package dao

import "errors"

var (
	URLNotFoundError = errors.New("URL not found")
)

type AbstractDAO interface {
	GetURLByShortLink(shortLink string) (string, error)
	AddShortURL(shortLink string, fullLink string) error
}