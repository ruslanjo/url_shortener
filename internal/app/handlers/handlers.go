package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ruslanjo/url_shortener/internal/config"
	"github.com/ruslanjo/url_shortener/internal/core"
)

var storage = make(map[string]string)

func Dispatcher(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		CreateShortUrl(w, req)
	case http.MethodGet:
		GetUrlByShortLink(w, req)
	default:
		http.Error(w, "Methods GET and POST allowed", http.StatusMethodNotAllowed)

	}
}

func CreateShortUrl(w http.ResponseWriter, req *http.Request) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	rawURL := string(data)
	rawURL = strings.Trim(rawURL, "\"")
	hash := core.HashString(rawURL)
	encodedURL := core.EncodeHash(hash)
	storage[encodedURL] = rawURL
	fmt.Println(storage)

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("http://%s/%s", config.Addr, encodedURL)))
}

func GetUrlByShortLink(w http.ResponseWriter, req *http.Request) {
	shortURL := req.URL.String()[1:]
	full, ok := storage[shortURL]
	if ok == false {
		http.Error(w, "Not found", http.StatusBadRequest)
		return
	}
	http.Redirect(w, req, full, http.StatusTemporaryRedirect)
}
