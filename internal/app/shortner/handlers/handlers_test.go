package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMultiplexer(t *testing.T) {
	jsonBody := []byte(`{"apple.com/store"}`)
	bodyReader := bytes.NewReader(jsonBody)
	requestURL := "localhost:8080/"

	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)

	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	Multiplexer(rec, req)

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
