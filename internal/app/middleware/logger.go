package middleware

import (
	"net/http"
	"time"

	"github.com/ruslanjo/url_shortener/internal/logger"
	"go.uber.org/zap"
)

/*
   Сведения о запросах должны содержать URI, метод запроса и время, затраченное на его выполнение.
   Сведения об ответах должны содержать код статуса и размер содержимого ответа.
*/

type RequestData struct {
	size       int
	statusCode int
}

type LogResponseWritter struct {
	http.ResponseWriter
	requestData RequestData
}

func (lw *LogResponseWritter) Write(b []byte) (int, error) {
	size, err := lw.ResponseWriter.Write(b)
	lw.requestData.size += size
	return size, err
}

func (lw *LogResponseWritter) WriteHeader(statusCode int) {
	lw.ResponseWriter.WriteHeader(statusCode)
	lw.requestData.statusCode = statusCode
}

func RequestLogger(h http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		reqData := RequestData{}
		lw := LogResponseWritter{
			ResponseWriter: w,
			requestData:    reqData,
		}
		h.ServeHTTP(&lw, r)

		duration := time.Since(start)
		logger.Log.Desugar().Info(
			"New request",
			zap.String("method", r.Method),
			zap.String("URI", r.URL.RequestURI()),
			zap.String("protocol", r.Proto),
			zap.String("referrer", r.Referer()),
			zap.Int("status code", lw.requestData.statusCode),
			zap.Int("response size", lw.requestData.size),
			zap.Duration("duration", duration),
		)
	}
	return http.HandlerFunc(fn)
}
