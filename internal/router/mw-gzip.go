package router

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

type compressWriter struct {
	http.ResponseWriter
}

func GzipMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		// decompress the incoming request
		isCompressed := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
		if isCompressed {
			zr, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Println("GzipMiddleware error at gzip.NewReader():", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(zr)
		}

		// compress the outgoing response
		acceptsGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		if acceptsGzip {
			w = &compressWriter{w}
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (cw *compressWriter) WriteHeader(statusCode int) {
	cw.ResponseWriter.Header().Set("Content-Encoding", "gzip")
	cw.ResponseWriter.WriteHeader(statusCode)
}

func (cw *compressWriter) Write(p []byte) (int, error) {
	gz := gzip.NewWriter(cw.ResponseWriter)
	defer gz.Close()
	return gz.Write(p)
}
