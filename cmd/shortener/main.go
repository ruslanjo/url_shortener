package main

import (
	"fmt"
	"net/http"

	"github.com/ruslanjo/url_shortener/internal/app/handlers"
)

func main() {
	host := "0.0.0.0"
	port := 8080
	addr := fmt.Sprintf("%s:%d", host, port)

	err := http.ListenAndServe(addr, http.HandlerFunc(handlers.Dispatcher))
	if err != nil {
		panic(err)
	}
}
