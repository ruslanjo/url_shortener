package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ruslanjo/url_shortener/internal/app/dao"
	"github.com/ruslanjo/url_shortener/internal/config"
	"github.com/ruslanjo/url_shortener/internal/core"
)

func Dispatcher(dao dao.AbstractDAO) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			CreateShortURLHandler(dao)(w, req)
		case http.MethodGet:
			GetURLByShortLinkHandler(dao)(w, req)
		default:
			http.Error(w, "Methods GET and POST allowed", http.StatusMethodNotAllowed)

		}
	}
}

func CreateShortURLHandler(dao dao.AbstractDAO) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data, err := io.ReadAll(req.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		rawURL := string(data)
		rawURL = strings.Trim(rawURL, "\"")
		encodedURL := core.GenerateShortURL(rawURL)

		err = dao.AddShortURL(encodedURL, rawURL)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("%s/%s", config.FullAddr, encodedURL)))
	}
}

func GetURLByShortLinkHandler(dao dao.AbstractDAO) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		shortURL := req.URL.String()[1:]
		full, err := dao.GetURLByShortLink(shortURL)
		if err != nil {
			http.Error(w, "Not found", http.StatusBadRequest)
			return
		}
		http.Redirect(w, req, full, http.StatusTemporaryRedirect)
	}
}
