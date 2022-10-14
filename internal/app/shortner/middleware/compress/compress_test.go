package compress

import (
	"net/http"
	"reflect"
	"testing"
)

func TestGzipHandler_GzipMiddleware(t *testing.T) {
	type args struct {
		next http.Handler
	}
	tests := []struct {
		name string
		g    GzipHandler
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := GzipHandler{}
			if got := g.GzipMiddleware(tt.args.next); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GzipHandler.GzipMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
