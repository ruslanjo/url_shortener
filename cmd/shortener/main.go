package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ruslanjo/url_shortener/internal/app/handlers"
	"github.com/ruslanjo/url_shortener/internal/app/middleware"
	"github.com/ruslanjo/url_shortener/internal/app/storage"
	"github.com/ruslanjo/url_shortener/internal/app/storage/disk"
	"github.com/ruslanjo/url_shortener/internal/config"
	"github.com/ruslanjo/url_shortener/internal/logger"
)

func setUpRouter(storage storage.Storage, db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger)
	r.Use(middleware.Compression)

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.CreateShortURLHandler(storage))
		r.Get("/{shortURL}", handlers.GetURLByShortLinkHandler(storage))
		r.Post("/api/shorten", handlers.GetShortURLJSONHandler(storage))
		r.Get("/ping", handlers.PingDB(db))
	})
	return r
}

func initStorage() (storage.Storage, *sql.DB) {
	if config.DSN == "" {
		urlDs := disk.NewURLDiskStorage(config.LocalStoragePath)
		storage := storage.NewHashMapStorage(urlDs)
		if err := storage.LoadFromDisk(); err != nil {
			log.Fatal(err)
		}
		logger.Log.Infoln("storage: memory and disk")
		return storage, nil
	}

	db := config.MustLoadDB()
	dbStorage := storage.NewPostgresStorage(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	storage.InitPostgres(ctx, db)
	logger.Log.Infoln("storage: Postgres")
	return &dbStorage, db
}

func main() {
	config.ConfigureApp()
	logger.Initialize("info")

	storage, dbDriver := initStorage()

	r := setUpRouter(storage, dbDriver)
	logger.Log.Infoln("Starting server")
	log.Fatal(http.ListenAndServe(config.ServerAddr, r))
}
