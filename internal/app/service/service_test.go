package service_test

import (
	"testing"

	"github.com/AndreyAD1/url-shortener/internal/app/service"
	"github.com/stretchr/testify/require"
)

func Test_getRandomString(t *testing.T) {
	tests := []struct {
		name string
		length int
	}{
		{"common", 10},
		{"empty", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			randomString := service.GetRandomString(tt.length)
			require.Equal(t, tt.length, len(randomString))
		})
	}
}
