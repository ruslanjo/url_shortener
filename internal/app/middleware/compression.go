package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/ruslanjo/url_shortener/internal/config"
)

type compressWriter struct {
	w     http.ResponseWriter
	cw    io.WriteCloser
	cType string
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:     w,
		cw:    gzip.NewWriter(w),
		cType: config.SelectedCompressionType,
	}
}

func (c *compressWriter) Write(data []byte) (int, error) {
	return c.cw.Write(data)
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) WriteHeader(statusCode int) {
	c.w.Header().Set("Content-Encoding", string(c.cType))
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.cw.Close()
}

type compressReader struct {
	body io.ReadCloser
	cr   io.ReadCloser
}

func newCompressReader(body io.ReadCloser) (*compressReader, error) {
	zipR, err := gzip.NewReader(body)
	if err != nil {
		return nil, err
	}
	return &compressReader{
		body: body,
		cr:   zipR,
	}, nil
}

func (c compressReader) Read(p []byte) (int, error) {
	return c.cr.Read(p)
}

func (c compressReader) Close() error {
	if err := c.body.Close(); err != nil {
		return err
	}
	return c.cr.Close()
}

func isApplicableContentType(contentType string) bool {
	applicable := [2]string{"application/json", "text/html"}
	contentType = strings.ToLower(contentType)
	for _, ct := range applicable {
		if contentType == ct {
			return true
		}
	}
	return false
}

func Compression(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		writer := w
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsCompr := strings.Contains(acceptEncoding, config.SelectedCompressionType)
		isApplicableCT := isApplicableContentType(r.Header.Get("Content-Type"))
		if supportsCompr && isApplicableCT{
			cw := newCompressWriter(w)
			writer = cw
			defer cw.Close()
		}

		isEncoded := r.Header.Get("Content-Encoding")
		sendsCompr := strings.Contains(isEncoded, config.SelectedCompressionType)
		if sendsCompr {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(writer, r)

	}
	return http.HandlerFunc(fn)
}
