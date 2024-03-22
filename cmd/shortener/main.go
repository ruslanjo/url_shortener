package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ruslanjo/url_shortener/internal/app/handlers"
	"github.com/ruslanjo/url_shortener/internal/app/middleware"
	"github.com/ruslanjo/url_shortener/internal/app/storage"
	"github.com/ruslanjo/url_shortener/internal/app/storage/disk"
	"github.com/ruslanjo/url_shortener/internal/config"
	"github.com/ruslanjo/url_shortener/internal/logger"
)

func setUpRouter(storage storage.Storage) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger)
	r.Use(middleware.Compression)

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.CreateShortURLHandler(storage))
		r.Get("/{shortURL}", handlers.GetURLByShortLinkHandler(storage))
		r.Post("/api/shorten", handlers.GetShortURLJSONHandler(storage))
		r.Post("/api/shorten/batch", handlers.BatchShortenHandler(storage))
		r.Get("/ping", handlers.PingDBHandler(storage))
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
	storage.InitPostgres(db)
	logger.Log.Infoln("storage: Postgres")
	return &dbStorage, db
}

func main() {
	config.ConfigureApp()
	logger.Initialize("info")

	storage, dbDriver := initStorage()

	r := setUpRouter(storage)
	logger.Log.Infoln("Starting server")
	log.Fatal(http.ListenAndServe(config.ServerAddr, r))
}
