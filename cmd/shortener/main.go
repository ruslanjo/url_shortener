package main

import (
	"net/http"

	"github.com/ruslanjo/url_shortener/internal/app/dao"
	"github.com/ruslanjo/url_shortener/internal/app/handlers"
	"github.com/ruslanjo/url_shortener/internal/config"
)

func main() {
	dao := &dao.HashMapDAO{}
	err := http.ListenAndServe(config.Addr, http.HandlerFunc(handlers.Dispatcher(dao)))
	if err != nil {
		panic(err)
	}
}
