package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/AndreyAD1/url-shortener/internal/app/service"
	"github.com/AndreyAD1/url-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortURLHandler(t *testing.T) {
	db := storage.NewStorage()
	URLService := service.Service{Storage: db}
	handler := ShortURLHandler(URLService)

	tests := []struct {
		Name                 string
		Method               string
		Path                 string
		Body                 string
		ExpectedResponseCode int
	}{
		{
			"invalid request method",
			http.MethodPut,
			"whatever",
			"",
			http.StatusMethodNotAllowed,
		},
		{
			"POST with invalid path",
			http.MethodPost,
			"localhost/some/path",
			"whatever",
			http.StatusNotFound,
		},
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

func TestShortURLHandler_GET(t *testing.T) {
	db := storage.NewStorage()
	URLService := service.Service{Storage: db}
	handler := ShortURLHandler(URLService)
	testURL := "http://test/url"
	parsedTestURL, err := url.Parse(testURL)
	require.NoError(t, err)
	testUrlID := "test-url-id"
	db.WriteURL(testUrlID, *parsedTestURL)

	tests := []struct {
		Name                 string
		Path                 string
		ExpectedResponseCode int
	}{
		{
			"Too long path",
			"http://localhost/some/path",
			http.StatusNotFound,
		},
		{
			"no URL id",
			"http://localhost/",
			http.StatusNotFound,
		},
		{
			"absent URL id",
			"http://localhost/absent-id",
			http.StatusBadRequest,
		},
		{
			"valid request",
			"http://localhost/" + testUrlID,
			http.StatusTemporaryRedirect,
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
			assert.Equal(t, tc.ExpectedResponseCode, result.StatusCode)
			if tc.ExpectedResponseCode == http.StatusTemporaryRedirect {
				require.Equal(t, testURL, result.Header.Get("Location"))
			}
		})
	}
}
