package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ruslanjo/url_shortener/internal/app/storage"
	"github.com/ruslanjo/url_shortener/internal/config"
	"github.com/ruslanjo/url_shortener/internal/core"
)

func CreateShortURLHandler(storage storage.Storage) http.HandlerFunc {
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

func GetURLByShortLinkHandler(storage storage.Storage) http.HandlerFunc {
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

func GetShortURLJSONHandler(storage storage.Storage) http.HandlerFunc {
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

		err = storage.AddShortURL(shortURL, input.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(output)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)

	}
}

func PingDB(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(fn)

}
