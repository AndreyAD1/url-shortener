package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	ZipWriter    io.Writer
	CommonWriter http.ResponseWriter
}

var compressedContentTypes = []string{
	"application/javascript",
	"application/json",
	"text/css",
	"text/html",
	"text/plain",
	"text/xml",
}

func (w gzipWriter) Write(b []byte) (int, error) {
	currentContentType := w.Header().Get("Content-Type")
	for _, compressedContentType := range compressedContentTypes {
		if strings.Contains(currentContentType, compressedContentType) {
			return w.ZipWriter.Write(b)
		}
	}
	return w.CommonWriter.Write(b)
}

func DecompressGzipRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") != "gzip" {
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

func CompressResponseToGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()
		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{w, gz, w}, r)
	})
}
