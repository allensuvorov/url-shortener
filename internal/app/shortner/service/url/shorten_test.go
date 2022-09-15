package url

import (
	"testing"
)

func TestShorten(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildHash(tt.args.s); got != tt.want {
				t.Errorf("Shorten() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getUniqShortHash(t *testing.T) {
	type args struct {
		h string
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			// TODO: Add test cases.
			name: "test shortening",
			args: args{"1234567890", "a/b/c"},
			want: "12345678",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getUniqShortHash(tt.args.h, tt.args.s); got != tt.want {
				t.Errorf("getUniqShortHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
