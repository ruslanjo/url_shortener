package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/ruslanjo/url_shortener/internal/app/storage"
	"github.com/ruslanjo/url_shortener/internal/app/storage/mocks"
	"github.com/ruslanjo/url_shortener/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ResponeInfo struct {
	code        int
	contentType string
	fullLink    string
	shortLink   string
}

type TestMeta []struct {
	name string
	info ResponeInfo
}

func getLink(resource string) string {
	return fmt.Sprintf("%s/%s", config.BaseServerReturnAddr, resource)
}

func TestCreateShortURLHandler(t *testing.T) {
	mockDao := &storage.HashMapStorage{}

	testSuits := []struct {
		name string
		info ResponeInfo
	}{
		{
			name: "Creation #1",
			info: ResponeInfo{
				code:        201,
				contentType: "text/plain",
				fullLink:    "https://ya.ru",
				shortLink:   "6YGS4ZUFRyR2pJ8QOIQoqw==",
			},
		},
	}

	for _, tt := range testSuits {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.info.fullLink))
			r.Header.Set("Content-Type", "text/plain")
			ctx := context.WithValue(r.Context(), config.CtxUserIDKey, "random")
			r = r.WithContext(ctx)
			CreateShortURLHandler(mockDao)(w, r)
			res := w.Result()

			assert.Equal(t, tt.info.code, res.StatusCode)

			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, getLink(tt.info.shortLink), string(body))
			assert.Equal(t, tt.info.contentType, res.Header.Get("Content-Type"))

		})
	}
}

func TestGetURLByShortLinkHandler(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	mockStorage.EXPECT().GetURLByShortLink("6YGS4ZUFRyR2pJ8QOIQoqw==").Return("https://ya.ru", nil)
	mockStorage.EXPECT().GetURLByShortLink(gomock.Any()).Return("", storage.ErrURLMappingNotFound)

	testSuits := TestMeta{
		{
			name: "Link exists",
			info: ResponeInfo{
				code:      307,
				fullLink:  "https://ya.ru",
				shortLink: "6YGS4ZUFRyR2pJ8QOIQoqw==",
			},
		},
		{
			name: "Link not exists",
			info: ResponeInfo{
				code:      400,
				fullLink:  "https://ufc.com",
				shortLink: "blablabla",
			},
		},
	}

	for _, tt := range testSuits {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tt.info.shortLink), nil)

			rtcx := chi.NewRouteContext()
			rtcx.URLParams.Add("shortURL", tt.info.shortLink)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rtcx))

			GetURLByShortLinkHandler(mockStorage)(w, r)

			res := w.Result()

			defer res.Body.Close()
			assert.Equal(t, tt.info.code, res.StatusCode)

			if tt.info.code == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.info.fullLink, res.Header.Get("Location"))
			}
		})
	}

}

func TestGetShortURLJSONHandler(t *testing.T) {
	storage := &storage.HashMapStorage{}
	url := "/api/shorten"
	testSuits := []struct {
		name         string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "empty body",
			body:         "",
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
		{
			name:         "wrong JSON",
			body:         `{"url": "bla}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
		{
			name:         "Success POST",
			body:         `{"url": "http://ya.ru"}`,
			expectedCode: http.StatusCreated,
			expectedBody: `{"result": "http://localhost:8080/G1VrRKTuc1JPsAnhGRj7Tw=="}`,
		},
	}

	for _, tt := range testSuits {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer([]byte(tt.body))
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, url, buf)
			ctx := context.WithValue(r.Context(), config.CtxUserIDKey, "random")
			r = r.WithContext(ctx)
			if len(tt.body) > 0 {
				r.Header.Set("Content-Type", "application/json")
			}
			GetShortURLJSONHandler(storage)(w, r)
			res := w.Result()
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf(err.Error())
			}
			assert.Equal(t, tt.expectedCode, res.StatusCode)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, string(body))
			}

		})
	}
}
