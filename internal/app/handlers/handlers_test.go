package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AndreyAD1/url-shortener/internal/app/service"
	"github.com/AndreyAD1/url-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateShortURLHandler(t *testing.T) {
	db := storage.NewStorage()
	URLService := service.Service{Storage: db}
	handler := CreateShortURLHandler(URLService)

	tests := []struct {
		Name                 string
		Method               string
		Path                 string
		Body                 string
		ExpectedResponseCode int
	}{
		{
			"POST with invalid URI",
			http.MethodPost,
			"http://localhost/",
			"invalidURI",
			http.StatusBadRequest,
		},
		{
			"POST happy pass",
			http.MethodPost,
			"http://localhost/",
			"http://some_test/url/test",
			http.StatusCreated,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			var err error
			var request *http.Request
			requestBody := bytes.NewBufferString(tc.Body)
			if tc.Body == "" {
				request, err = http.NewRequest(tc.Method, tc.Path, nil)
				require.NoError(t, err)
			} else {
				request, err = http.NewRequest(tc.Method, tc.Path, requestBody)
				require.NoError(t, err)
			}
			recorder := httptest.NewRecorder()
			handler(recorder, request)
			result := recorder.Result()
			defer result.Body.Close()

			assert.Equal(t, tc.ExpectedResponseCode, result.StatusCode)
		})
	}
}

func TestGetFullURLHandler(t *testing.T) {
	db := storage.NewStorage()
	URLService := service.Service{Storage: db}
	handler := GetFullURLHandler(URLService)

	tests := []struct {
		Name                 string
		Path                 string
		ExpectedResponseCode int
	}{
		{
			"absent URL id",
			"http://localhost/absent-id",
			http.StatusBadRequest,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, tc.Path, nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			handler(recorder, request)
			result := recorder.Result()
			defer result.Body.Close()
			require.Equal(t, tc.ExpectedResponseCode, result.StatusCode)
		})
	}
}
