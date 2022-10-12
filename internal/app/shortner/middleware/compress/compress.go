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

// Reader
type gzipReader struct {
	rc  io.ReadCloser
	gzr *gzip.Reader
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

// middleware принимает параметром Handler и возвращает тоже Handler.
func (g GzipHandler) GzipMiddleware(next http.Handler) http.Handler {
	// собираем Handler приведением типа
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Compress/GzipMiddleware: Hi, I'm GzipMiddleware ")

		// check if request is gzip-encoded
		// if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		// 	log.Println("Compress/Middleware: POST request Content-Encoding == gzip")
		// 	// TODO: read response

		// 	gz, _ := gzip.NewReader(r.Body)
		// 	gzr := gzipReader{
		// 		rc:  r.Body,
		// 		gzr: gz,
		// 	}
		// 	r.Body = gzr
		// 	defer gz.Close()

		// }

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
			// next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
		} else {

			log.Println("Compress/Middleware: gzip not accepted, header:", r.Header["Accept-Encoding"])

			// замыкание — используем ServeHTTP следующего хендлера
			next.ServeHTTP(w, r)
		}

		log.Println("Compress/GzipMiddleware: Bye! ")
	})
}

// type WrappedResponseWriter struct {
// 	http.ResponseWriter
// 	gw *gzip.Writer
// }

// // middleware принимает параметром Handler и возвращает тоже Handler.
// func (g GzipHandler) GzipMiddleware(next http.Handler) http.Handler {
// 	// собираем Handler приведением типа

// 		// w.Header().Set("Access-Control-Allow-Origin", "*")

// 		// if r.Method == http.MethodPost {
// 		// 	log.Println("Handler/Middleware: request method = post ")

// 		// 	// переменная rc будет равна r.Body или *gzip.Reader
// 		// 	var rc io.Reader
// 		// 	var rc io.ReadCloser

// 		// if r.Header.Get(`Content-Encoding`) == `gzip` {
// 		// 		log.Println("Handler/Middleware: POST request Content-Encoding == gzip")
// 		// 		gz, err := gzip.NewReader(r.Body)
// 		// 		if err != nil {
// 		// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		// 			return
// 		// 		}
// 		// 		rc = gz
// 		// 		defer gz.Close()

// 		// 	} else {
// 		// 		log.Println("Handler/Middleware: POST request Content-Encoding is not gzip")
// 		// 		rc = r.Body
// 		// 		defer r.Body.Close()
// 		// 	}
// 		// 	r.Body = rc
// 		// 	b, err := io.ReadAll(r.Body)
// 		// 	log.Println("Handler/Middleware: POST request body:", string(b))

// 		// 	if err != nil {
// 		// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		// 		return
// 		// 	}
// 		// 	log.Printf("Handler/Middleware: Length: %d", len(b))

// 		// 	next.ServeHTTP(w, r)
// 		// }

// 		// проверяем, что клиент поддерживает gzip-сжатие
// 		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
// 			log.Println("Handler/Middleware: gzip not accepted, header:", r.Header["Accept-Encoding"])
// 			// если gzip не поддерживается, передаём управление
// 			// дальше без изменений
// 			next.ServeHTTP(w, r)
// 			return
// 		}
// 		log.Println("Handler/Middleware: gzip is accepted, header:", r.Header["Accept-Encoding"])
// 		// создаём gzip.Writer поверх текущего w
// 		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
// 		if err != nil {
// 			io.WriteString(w, err.Error())
// 			return
// 		}
// 		defer gz.Close()

// 		w.Header().Set("Content-Encoding", "gzip")
// 		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
// 		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)

// 		// замыкание — используем ServeHTTP следующего хендлера

// 		log.Println("Handler/Middleware: Bye! ")

// 	})
// }
