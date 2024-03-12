package storage

import (
	"context"
	"database/sql"

	"github.com/ruslanjo/url_shortener/internal/app/storage/models"
)

func NewPostgresStorage(db *sql.DB) postgresStorage {
	pg := postgresStorage{db}
	return pg
}

type postgresStorage struct {
	db *sql.DB
}

func (pg *postgresStorage) GetURLByShortLink(shortLink string) (string, error) {
	row := pg.db.QueryRow(
		"select url from urls where alias=$1", shortLink,
	)

	url := new(string)
	if err := row.Scan(url); err != nil {
		return "", err
	}
	return *url, nil
}

func (pg *postgresStorage) AddShortURL(shortLink string, fullLink string) error {
	_, err := pg.db.Exec(
		"insert into urls(url, alias) values ($1, $2)", fullLink, shortLink,
	)
	return err
}

func (pg *postgresStorage) SaveURLBatched(ctx context.Context, data []models.URLBatch) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}

	for _, ent := range data {
		_, err := tx.ExecContext(ctx,
			"insert into urls(url, alias) values ($1, $2)",
			ent.OriginalURL, ent.ShortURL)
		if err != nil {
			tx.Rollback()
			return err
		}

	}
	tx.Commit()
	return nil
}
