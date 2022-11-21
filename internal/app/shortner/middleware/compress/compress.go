package compress

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

type gzipReader struct {
	rc  io.ReadCloser
	gzr *gzip.Reader
}

func (g gzipReader) Read(b []byte) (n int, err error) {
	return g.gzr.Read(b)
}

func (g gzipReader) Close() error {
	err := g.rc.Close()
	if err != nil {
		return fmt.Errorf("failed decompress data: %v", err)
	}
	return g.gzr.Close()
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Compress/GzipMiddleware: Hi, I'm GzipMiddleware ")

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			log.Println("Compress/GzipMiddleware: POST request Content-Encoding == gzip")

			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Fatal("Compress/GzipMiddleware: error - could not create gzip reader: ", err)
			}

			gzr := gzipReader{
				rc:  r.Body,
				gzr: gz,
			}
			defer gzr.Close()
			r.Body = gzr
		}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			log.Println("Compress/Middleware: gzip can be accepted, header:", r.Header["Accept-Encoding"])

			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			defer gz.Close()

			w.Header().Set("Content-Encoding", "gzip")
			w = gzipWriter{
				ResponseWriter: w,
				Writer:         gz,
			}
			log.Println("Compress/Middleware: gzip.Writer - end")
		}

		log.Println("Compress/Middleware: w and r ready to pass to next")
		next.ServeHTTP(w, r)
		log.Println("Compress/GzipMiddleware: Bye! ")
	})
}
