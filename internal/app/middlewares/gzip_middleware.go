package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(`Content-Encoding`) != `gzip` {
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		unzippedBody, err := io.ReadAll(gz)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		unzippedReader := bytes.NewBuffer(unzippedBody)
		unzippedRequest, err := http.NewRequest(
			r.Method,
			r.RequestURI,
			unzippedReader,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		unzippedRequest.Header = r.Header
		next.ServeHTTP(w, unzippedRequest)
	})
}
