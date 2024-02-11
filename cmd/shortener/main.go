package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ruslanjo/url_shortener/internal/app/dao"
	"github.com/ruslanjo/url_shortener/internal/app/handlers"
	"github.com/ruslanjo/url_shortener/internal/config"
)

func setUpRouter(dao dao.AbstractDAO) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.CreateShortURLHandler(dao))
		r.Get("/{shortURL}", handlers.GetURLByShortLinkHandler(dao))
	})
	return r
}

func main() {

	dao := &dao.HashMapDAO{}
	r := setUpRouter(dao)
	log.Fatal(http.ListenAndServe(config.Addr, r))
}
