package main

import (
	"net/http"

	"github.com/ruslanjo/url_shortener/internal/app/handlers"
	"github.com/ruslanjo/url_shortener/internal/config"
)

func main() {
	err := http.ListenAndServe(config.Addr, http.HandlerFunc(handlers.Dispatcher))
	if err != nil {
		panic(err)
	}
}
