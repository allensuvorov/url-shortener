package compress

import (
	"log"
	"net/http"
	"strings"
)

type GzipHandler struct {
}

// middleware принимает параметром Handler и возвращает тоже Handler.
func (g GzipHandler) GzipMiddleware(next http.Handler) http.Handler {
	// собираем Handler приведением типа
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handler/GzipMiddleware: Hi, I'm GzipMiddleware ")

		// check if request is gzip-encoded
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			log.Println("Handler/Middleware: POST request Content-Encoding == gzip")
			// TODO: read response
		}

		// check if sender can accept is gzip-encoded response
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			log.Println("Handler/Middleware: gzip can be accepted, header:", r.Header["Accept-Encoding"])
			// TODO: create gziped response

		}

		// log.Println("Handler/Middleware: gzip not accepted, header:", r.Header["Accept-Encoding"])
		// если gzip не поддерживается, передаём управление
		// дальше без изменений
		next.ServeHTTP(w, r)
	})
}

// type WrappedResponseWriter struct {
// 	http.ResponseWriter
// 	gw *gzip.Writer
// }

// type gzipWriter struct {
// 	http.ResponseWriter
// 	Writer io.Writer
// }

// func (w gzipWriter) Write(b []byte) (int, error) {
// 	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
// 	return w.Writer.Write(b)
// }

// // Reader
// type gzipReader struct {
// 	rc  io.ReadCloser
// 	gzr *gzip.Reader
// }

// func (g gzipReader) Read(p []byte) (n int, err error) {}
// func (g gzipReader) Close() error {}

// func main() {
// 	req, _ := http.NewRequest(http.MethodGet, "", nil)
// 	gz, _ := gzip.NewReader(req.Body)
// 	gzr := gzipReader{
// 		rc:  req.Body,
// 		gzr: gz,
// 	}
// 	req.Body = gzr
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
