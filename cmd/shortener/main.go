package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	h "github.com/ruslanjo/url_shortener/internal/app/handlers"
	mw "github.com/ruslanjo/url_shortener/internal/app/middleware"
	"github.com/ruslanjo/url_shortener/internal/app/storage"
	"github.com/ruslanjo/url_shortener/internal/app/storage/disk"
	"github.com/ruslanjo/url_shortener/internal/config"
	"github.com/ruslanjo/url_shortener/internal/logger"
)

func setUpRouter(storage storage.Storage, tokenGen mw.TokenGenerator) *chi.Mux {
	r := chi.NewRouter()
	r.Use(mw.RequestLogger)
	r.Use(mw.Compression)

	r.Post("/", mw.Signup(h.CreateShortURLHandler(storage), tokenGen))
	r.Get("/{shortURL}", h.GetURLByShortLinkHandler(storage))
	r.Get("/ping", h.PingDBHandler(storage))
	r.Route("/api/shorten", func(r chi.Router) {
		r.Post("/", mw.Signup(h.GetShortURLJSONHandler(storage), tokenGen))
		r.Post("/batch", mw.Signup(h.BatchShortenHandler(storage), tokenGen))
	})
	r.Route("/api/user/urls", func(r chi.Router) {
		r.Get("/", h.GetUserURLsHandler(storage, tokenGen))
		r.Delete("/", mw.Signup(h.DeleteURLsHandler(storage), tokenGen))
	})
	return r
}

func initStorage() storage.Storage {
	if config.DSN == "" {
		urlDs := disk.NewURLDiskStorage(config.LocalStoragePath)
		storage := storage.NewHashMapStorage(urlDs)
		if err := storage.LoadFromDisk(); err != nil {
			log.Fatal(err)
		}
		logger.Log.Infoln("storage: memory and disk")
		return storage
	}

	dbDriver := config.MustLoadDB()
	dbStorage := storage.NewPostgresStorage(dbDriver)
	storage.InitPostgres(dbDriver)
	logger.Log.Infoln("storage: Postgres")
	return &dbStorage
}

func main() {
	config.ConfigureApp()
	logger.Initialize("info")

	storage := initStorage()

	tokenGenerator := mw.TokenGenerator{}

	r := setUpRouter(storage, tokenGenerator)
	logger.Log.Infoln("Starting server")
	log.Fatal(http.ListenAndServe(config.ServerAddr, r))
}
