package storage

import (
	"database/sql"
	"fmt"
	"log"
)

func InitPostgres(db *sql.DB) {
	q := `
	CREATE TABLE IF NOT EXISTS urls(
		id serial CONSTRAINT urls_pl PRIMARY KEY,
		url varchar(512),
		alias varchar(512),
		uuid varchar(64),
		CONSTRAINT url_uniq UNIQUE (url),
		CONSTRAINT alias_uniq UNIQUE (alias)
	);
	`
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(
			fmt.Errorf("error while creating transaction: %w", err),
		)
	}

	if _, err := tx.Exec(q); err != nil {
		log.Fatal(
			fmt.Errorf("error while creating DB table: %w", err),
		)
		tx.Rollback()
	}
	tx.Commit()

}
