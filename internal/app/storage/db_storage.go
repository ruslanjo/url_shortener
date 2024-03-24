package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

func (pg *postgresStorage) AddShortURL(
	shortLink string,
	fullLink string,
	UUID string,
) error {
	_, err := pg.db.Exec(
		"insert into urls(url, alias, uuid) values ($1, $2, $3)", fullLink, shortLink, UUID,
	)
	if err != nil {
		err = handleConstraintViolation(err)

	}

	return err
}

func (pg *postgresStorage) SaveURLBatched(
	ctx context.Context,
	data []models.URLBatch,
	UUID string,
) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}

	for _, ent := range data {
		_, err := tx.ExecContext(ctx,
			"insert into urls(url, alias, uuid) values ($1, $2, $3)",
			ent.OriginalURL, ent.ShortURL, UUID)
		if err != nil {
			tx.Rollback()
			err = handleConstraintViolation(err)
			return err
		}

	}
	tx.Commit()
	return nil
}

func (pg *postgresStorage) GetUserURLs(UUID string) ([]models.URL, error) {
	var urls []models.URL

	rows, err := pg.db.Query(
		"select url, alias from urls where uuid = $1", UUID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u models.URL
		err = rows.Scan(&u.OriginalURL, &u.ShortURL)
		if err != nil {
			return nil, err
		}

		urls = append(urls, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return urls, nil

}

func (pg *postgresStorage) PingContext(ctx context.Context) error {
	return pg.db.PingContext(ctx)
}

func handleConstraintViolation(err error) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		err = ErrIntegityViolation
	}
	return err
}
