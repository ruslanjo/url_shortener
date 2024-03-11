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
		CONSTRAINT url_uniq UNIQUE (url),
		CONSTRAINT alias_uniq UNIQUE (alias)
	);
	`
	if _, err := db.Exec(q); err != nil {
		log.Fatal(
			fmt.Errorf("error while creating DB table: %w", err),
		)
	}

}
