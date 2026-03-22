package middlewares

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// If not compressed by gzip, function compresses the plain text

func CompressionHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

		encoding := r.Header.Get("Accept-Encoding")

		if !strings.Contains(encoding, "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		w = &gzipWriter{ResponseWriter: w, Writer: gz}
		next.ServeHTTP(w, r)
	})
}

type gzipWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (gz *gzipWriter) Write (b []byte) (int, error) {
	return gz.Writer.Write(b);
}