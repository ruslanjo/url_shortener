package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ruslanjo/url_shortener/internal/app/storage"
	"github.com/ruslanjo/url_shortener/internal/config"
	"github.com/ruslanjo/url_shortener/internal/core"
)

func CreateShortURLHandler(storage storage.AbstractStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data, err := io.ReadAll(req.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		rawURL := string(data)
		encodedURL := core.GenerateShortURL(rawURL)

		err = storage.AddShortURL(encodedURL, rawURL)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("%s/%s", config.BaseServerReturnAddr, encodedURL)))
	}
}

func GetURLByShortLinkHandler(storage storage.AbstractStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		shortURL := chi.URLParam(req, "shortURL")
		full, err := storage.GetURLByShortLink(shortURL)
		if err != nil {
			http.Error(w, "Not found", http.StatusBadRequest)
			return
		}
		http.Redirect(w, req, full, http.StatusTemporaryRedirect)
	}
}

func GetShortURLJSONHandler(storage storage.AbstractStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			input struct {
				URL string `json:"url"`
			}
			output struct {
				URL string `json:"result"`
			}
		)

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if len(input.URL) == 0 {
			http.Error(w, "Please, pass url with length gt 0", http.StatusBadRequest)
			return
		}
		shortURL := core.GenerateShortURL(input.URL)
		output.URL = fmt.Sprintf(
			"%s/%s",
			config.BaseServerReturnAddr, shortURL,
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		response, err := json.Marshal(output)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = storage.AddShortURL(shortURL, input.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(response)

	}
}
