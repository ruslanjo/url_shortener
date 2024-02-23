package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ruslanjo/url_shortener/internal/app/storage"
	"github.com/ruslanjo/url_shortener/internal/app/handlers"
	"github.com/ruslanjo/url_shortener/internal/app/middleware"
	"github.com/ruslanjo/url_shortener/internal/config"
	"github.com/ruslanjo/url_shortener/internal/logger"
)

func setUpRouter(storage storage.AbstractStorage) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestLogger)

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.CreateShortURLHandler(storage))
		r.Get("/{shortURL}", handlers.GetURLByShortLinkHandler(storage))
	})
	return r
}

func main() {
	config.ConfigureApp()
	logger.Initialize("info")

	dao := &storage.HashMapStorage{}
	r := setUpRouter(dao)
	logger.Log.Infoln("Starting server")
	log.Fatal(http.ListenAndServe(config.ServerAddr, r))
}
