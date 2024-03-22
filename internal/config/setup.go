package config

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func parseFlags() {
	flag.StringVar(&ServerAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&Protocol, "p", "http", "Protocol to run server")
	flag.StringVar(&BaseServerReturnAddr, "b", "http://localhost:8080", "Base addres of URL shortener")
	flag.StringVar(&LocalStoragePath, "f", "/tmp/short-url-db.json", "Local DB addr")
	flag.StringVar(&DSN, "d", "", "Database connection string")
	flag.Parse()
}

func setEnvsToConfig(envName string, variable *string) {
	if envValue, ok := os.LookupEnv(envName); ok {
		*variable = envValue
	}
}

func ConfigureApp() {
	parseFlags()

	envs := map[string]*string{
		"RUN_ADDR":          &ServerAddr,
		"BASE_URL":          &BaseServerReturnAddr,
		"FILE_STORAGE_PATH": &LocalStoragePath,
		"DATABASE_DSN":      &DSN,
	}

	for k, v := range envs {
		setEnvsToConfig(k, v)
	}
}

func MustLoadDB() *sql.DB {
	db, err := sql.Open("pgx", DSN)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}
	return db
}
