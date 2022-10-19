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
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

// gzipReader is a Readcloser. As is r.body
type gzipReader struct {
	rc  io.ReadCloser // here we will put original request.Body
	gzr *gzip.Reader  // here we will put new reader
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

// GzipMiddleware принимает параметром Handler и возвращает тоже Handler.
func GzipMiddleware(next http.Handler) http.Handler {
	// собираем Handler приведением типа
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Compress/GzipMiddleware: Hi, I'm GzipMiddleware ")

		// check if request is gzip-encoded
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			log.Println("Compress/Middleware: POST request Content-Encoding == gzip")

			// создаём *gzip.Reader, который будет читать тело запроса
			// и распаковывать его
			gz, err := gzip.NewReader(r.Body) // *Reader - type Reader struct {Header}
			if err != nil {
				log.Fatal(err)
			}

			gzr := gzipReader{ // new ReadCloser
				rc:  r.Body, // a ReadCloser - gets the original request.Body
				gzr: gz,     // *gzip.Reader of r.Body
			}
			defer gzr.Close()
			r.Body = gzr // struct that is ReadCloser
		}

		// check if sender can accept gzip-encoded response
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			log.Println("Compress/Middleware: gzip can be accepted, header:", r.Header["Accept-Encoding"])

			// создаём gzip.Writer поверх текущего w
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			defer gz.Close()

			w.Header().Set("Content-Encoding", "gzip")
			// передаём обработчику страницы переменную типа gzipWriter для вывода данных
			w = gzipWriter{
				ResponseWriter: w,
				Writer:         gz,
			}
			log.Println("Compress/Middleware: gzip.Writer - end")
		}

		// замыкание — используем ServeHTTP следующего хендлера
		log.Println("Compress/Middleware: w and r ready to pass to next")
		next.ServeHTTP(w, r)
		log.Println("Compress/GzipMiddleware: Bye! ")
	})
}
