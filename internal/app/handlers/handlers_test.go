package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ruslanjo/url_shortener/internal/app/dao"
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
	return fmt.Sprintf("%s/%s", config.FullAddr, resource)
}

func TestCreateShortURLHandler(t *testing.T) {
	mockDao := &dao.HashMapDAO{}

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

	mockDao := &dao.HashMapDAO{}
	mockDao.InitStorage(
		map[string]string{"6YGS4ZUFRyR2pJ8QOIQoqw==": "https://ya.ru"},
	)

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

			GetURLByShortLinkHandler(mockDao)(w, r)

			res := w.Result()
			assert.Equal(t, tt.info.code, res.StatusCode)

			if tt.info.code == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.info.fullLink, res.Header.Get("Location"))
			}
		})
	}

}
