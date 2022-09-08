package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostURL(t *testing.T) {

	jsonBody := []byte("http://www.apple.com/store")
	bodyReader := bytes.NewReader(jsonBody)
	requestURL := "localhost:8080/"

	req, err := http.NewRequest("POST", requestURL, bodyReader)

	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	PostURL(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status Created; got %v", res.Status)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("coult not read respons: %v", err)
	}
	if string(b) != "http://localhost:8080/a7d59904" {
		t.Fatalf("expected http://localhost:8080/a7d59904; got %s", string(b))
	}
}

func TestGetURL(t *testing.T) {
	// TODO: create a record for the the URL, to that this test can work independatly
	requestURL := "http://localhost:8080/a7d59904"
	req, err := http.NewRequest("GET", requestURL, nil)

	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	GetURL(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("expected status TemporaryRedirect; got %v", res.Status)
	}

	if res.Header.Get("Location") != "http://www.apple.com/store" {
		t.Errorf("Expected Header http://www.apple.com/store, got %s", res.Header.Get("Location"))
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
