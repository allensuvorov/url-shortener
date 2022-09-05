package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostHandler(t *testing.T) {
	// type post struct {
	// 	method string
	// 	URL    string
	// 	body   io.Reader
	// 	rbody  string
	// 	status int
	// 	err    string
	// }

	// tt := []struct {
	// 		name: "post req-res",
	// 			method: "POST",
	// 			URL: "localhost:8080/",
	// 			body: bytes.NewReader([]byte(`{"apple.com/store"}`)),
	// 			rbody: "http://localhost:8080/7a999481",
	// 			status:
	// }

	jsonBody := []byte(`{"apple.com/store"}`)
	bodyReader := bytes.NewReader(jsonBody)
	requestURL := "localhost:8080/"

	req, err := http.NewRequest("POST", requestURL, bodyReader)

	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	PostHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status Created; got %v", res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("coult not read respons: %v", err)
	}
	if string(b) != "http://localhost:8080/7a999481" {
		t.Fatalf("expected http://localhost:8080/7a999481; got %s", string(b))
	}
}

func TestGetHandeler(t *testing.T) {
	// type get struct {
	// 	method string
	// 	URL    string
	// 	header string // location
	// 	status int
	// 	err    string
	// }

	requestURL := "http://localhost:8080/7a999481"
	req, err := http.NewRequest("GET", requestURL, nil)

	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	GetHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("expected status TemporaryRedirect; got %v", res.Status)
	}

	if res.Header.Get("Location") != `{"apple.com/store"}` {
		t.Errorf("Expected Header apple.com/store, got %s", res.Header.Get("Location"))
	}

}

// func TestMultiplexer(t *testing.T) {
// 	type args struct {
// 		w http.ResponseWriter
// 		r *http.Request
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			Multiplexer(tt.args.w, tt.args.r)
// 		})
// 	}
// }