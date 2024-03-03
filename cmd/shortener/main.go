package main

import (
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

func setUpRouter(storage storage.AbstractStorage) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger)
	r.Use(middleware.Compression)

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.CreateShortURLHandler(storage))
		r.Get("/{shortURL}", handlers.GetURLByShortLinkHandler(storage))
		r.Post("/api/shorten", handlers.GetShortURLJSONHandler(storage))
	})
	return r
}

func main() {
	config.ConfigureApp()
	logger.Initialize("info")

	ds := disk.DiskStorage{Path: "storage.txt"}
	url_ds := disk.NewUrlDiskStorage(ds)
	storage := storage.NewHashMapStorage(url_ds)
	if err := storage.LoadFromDisk(); err != nil {
		log.Fatal(err)
	}
	r := setUpRouter(storage)
	logger.Log.Infoln("Starting server")
	log.Fatal(http.ListenAndServe(config.ServerAddr, r))
}
