package service_test

import (
	"testing"

	"github.com/AndreyAD1/url-shortener/internal/app/service"
)

func Test_getRandomString(t *testing.T) {
	type args struct {
		n int
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
			if got := service.GetRandomString(tt.args.n); got != tt.want {
				t.Errorf("getRandomString() = %v, want %v", got, tt.want)
			}
		})
	}
}
