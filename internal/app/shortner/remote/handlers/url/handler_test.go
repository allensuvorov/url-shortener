package url

import (
	"net/http"
	"testing"
)

func TestURLHandler_Create(t *testing.T) {
	type fields struct {
		urlService URLService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uh := URLHandler{
				urlService: tt.fields.urlService,
			}
			uh.Create(tt.args.w, tt.args.r)
		})
	}
}

func TestURLHandler_Get(t *testing.T) {
	type fields struct {
		urlService URLService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uh := URLHandler{
				urlService: tt.fields.urlService,
			}
			uh.Get(tt.args.w, tt.args.r)
		})
	}
}
