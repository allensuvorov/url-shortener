package compress

import (
	"compress/gzip"
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
	gzr *gzip.Reader  // here we will put new reader -
}

func (g gzipReader) Read(b []byte) (n int, err error) {
	return g.gzr.Read(b)
}

func (g gzipReader) Close() error {
	g.rc.Close()
	g.gzr.Close()
	return nil
}

type GzipHandler struct {
}

/* Logic of GzipMiddleware

- We are trying to give our server ability to decompress and compress data exchange.
- To do that we will add a middleware - GzipMiddleware - between router and handlers.
- GzipMiddleware "wraps" our handlers (“оборачивает” обработчики) to add gzip encoding to the flow.

- To do that GzipMiddleware takes:
	- our handler with it's arguments:
		- ResponseWriter - type interface - any objects with a required set of functionality.
		- pointer to a Request struct, with fields including:
			- url - path we can route to call the right handler
			- header - the list of paramenter, including encoding headers.
			- body - type interface io.ReadCloser ()
- Body carries the primary payload of the request, and we are going to decompress the payload.
- Body is any object with ReadCloser functionality. ReadCloser is the interface that groups the basic Read and Close methods.
- To decompress the data in body, we need to
	- read the data,
	- decompress this data,
	- create new body,
	- pass decompressed data to the new body,
	- close body reader.

- To do that we will build an object - a struct, that has read and close methods.

- returns a new handler function that calls our handler, with it's method ServeHTTP.
- In that call it passes updated arguments to that handler
*/

// GzipMiddleware принимает параметром Handler и возвращает тоже Handler.
func (g GzipHandler) GzipMiddleware(next http.Handler) http.Handler {
	// собираем Handler приведением типа
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Compress/GzipMiddleware: Hi, I'm GzipMiddleware ")

		// check if request is gzip-encoded
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			log.Println("Compress/Middleware: POST request Content-Encoding == gzip")

			// создаём *gzip.Reader, который будет читать тело запроса
			// и распаковывать его
			gz, _ := gzip.NewReader(r.Body) // *Reader - type Reader struct {Header}
			gzr := gzipReader{              // new ReadCloser
				rc:  r.Body, // a ReadCloser - gets the original request.Body
				gzr: gz,     // *gzip.Reader of r.Body
			}
			defer gzr.Close()
			r.Body = gzr // struct that is ReadCloser

		}

		// check if sender can accept is gzip-encoded response
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			log.Println("Compress/Middleware: gzip can be accepted, header:", r.Header["Accept-Encoding"])
			// TODO: create gziped response

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

		} else {

			log.Println("Compress/Middleware: gzip not accepted, header:", r.Header["Accept-Encoding"])

		}
		// замыкание — используем ServeHTTP следующего хендлера
		log.Println("Compress/Middleware: w and r ready to pass to next")

		next.ServeHTTP(w, r)

		log.Println("Compress/GzipMiddleware: Bye! ")
	})
}