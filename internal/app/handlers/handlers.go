package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ruslanjo/url_shortener/internal/app/middleware"
	"github.com/ruslanjo/url_shortener/internal/app/storage"
	"github.com/ruslanjo/url_shortener/internal/app/storage/models"
	"github.com/ruslanjo/url_shortener/internal/config"
	"github.com/ruslanjo/url_shortener/internal/core"
)

func getUserIDFromCtx(ctx context.Context) (string, error) {
	userIDCtx := ctx.Value(config.CtxUserIDKey)
	userID, ok := userIDCtx.(string)
	if !ok {
		return "", fmt.Errorf("error while receiving request's userID")
	}
	return userID, nil
}

func CreateShortURLHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var statusCode int = http.StatusCreated

		userID, err := getUserIDFromCtx(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rawURL := string(data)
		encodedURL := core.GenerateShortURL(rawURL)

		err = store.AddShortURL(encodedURL, rawURL, userID)
		if err != nil && !errors.Is(err, storage.ErrIntegityViolation) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, storage.ErrIntegityViolation) {
			statusCode = http.StatusConflict
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(statusCode)
		w.Write([]byte(fmt.Sprintf("%s/%s", config.BaseServerReturnAddr, encodedURL)))
	}
}

func GetURLByShortLinkHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		shortURL := chi.URLParam(r, "shortURL")
		full, err := store.GetURLByShortLink(shortURL)
		if err != nil {
			if errors.Is(err, storage.ErrEntityDeleted) {
				http.Error(w, err.Error(), http.StatusGone)
				return
			}
			http.Error(w, "Not found", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, full, http.StatusTemporaryRedirect)
	}
}

func GetShortURLJSONHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			input struct {
				URL string `json:"url"`
			}
			output struct {
				URL string `json:"result"`
			}
			statusCode int = http.StatusCreated
		)
		userID, err := getUserIDFromCtx(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewDecoder(r.Body).Decode(&input)
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

		err = store.AddShortURL(shortURL, input.URL, userID)
		if err != nil && !errors.Is(err, storage.ErrIntegityViolation) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if errors.Is(err, storage.ErrIntegityViolation) {
			statusCode = http.StatusConflict
		}
		response, err := json.Marshal(output)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write(response)

	}
}

func PingDBHandler(store storage.Storage) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := store.PingContext(ctx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(fn)
}

func BatchShortenHandler(store storage.Storage) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var inputBatch []models.URLBatch

		userID, err := getUserIDFromCtx(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = json.NewDecoder(r.Body).Decode(&inputBatch); err != nil {
			http.Error(w, "could not decode body to JSON", http.StatusBadRequest)
			return
		}
		if len(inputBatch) == 0 {
			http.Error(w, "passed empty JSON body", http.StatusBadRequest)
			return
		}

		for i := range inputBatch {
			inputBatch[i].ShortURL = core.GenerateShortURL(inputBatch[i].OriginalURL)
		}

		err = store.SaveURLBatched(r.Context(), inputBatch, userID)

		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("error while saving data: %v", err),
				http.StatusBadRequest,
			)
			return
		}

		for i := range inputBatch {
			inputBatch[i].ShortURL = fmt.Sprintf(
				"%s/%s",
				config.BaseServerReturnAddr, inputBatch[i].ShortURL,
			)
		}

		response, err := json.Marshal(inputBatch)
		if err != nil {
			http.Error(w, fmt.Sprintf("error while marshaling response: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	}
	return http.HandlerFunc(fn)
}

func GetUserURLsHandler(
	store storage.Storage,
	tokenGen middleware.TokenGenerator,
) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(config.AuthCookie)

		if err == http.ErrNoCookie {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := tokenGen.GetClaims(cookie.Value)
		if err != nil {
			var statusCode int
			switch {
			case errors.Is(err, middleware.ErrTokenNotValid):
				statusCode = http.StatusUnauthorized
			default:
				statusCode = http.StatusInternalServerError
			}
			http.Error(w, err.Error(), statusCode)
			return
		}

		urls, err := store.GetUserURLs(claims.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(urls) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		for i := range urls {
			urls[i].ShortURL = fmt.Sprintf(
				"%s/%s",
				config.BaseServerReturnAddr, urls[i].ShortURL,
			)
		}

		response, err := json.Marshal(urls)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
	return fn
}

func DeleteURLsHandler(store storage.Storage) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var shortURLs []string

		userID, err := getUserIDFromCtx(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewDecoder(r.Body).Decode(&shortURLs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		go store.DeleteURLs(context.Background(), shortURLs, userID)
		w.WriteHeader(http.StatusAccepted)
	}
	return fn
}
