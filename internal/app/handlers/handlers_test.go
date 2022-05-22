package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/AndreyAD1/url-shortener/internal/app/service"
	"github.com/AndreyAD1/url-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateShortURLHandler(t *testing.T) {
	db, err := storage.NewStorage("")
	require.NoError(t, err)
	URLService := service.Service{Storage: db}
	handlerContainer := HandlerContainer{URLService: URLService}

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
			handlerContainer.CreateShortURLHandler()(recorder, request)
			result := recorder.Result()
			defer result.Body.Close()

			assert.Equal(t, tc.ExpectedResponseCode, result.StatusCode)
		})
	}
}

func TestGetFullURLHandler(t *testing.T) {
	db, err := storage.NewStorage("")
	require.NoError(t, err)
	URLService := service.Service{Storage: db}
	handlerContainer := HandlerContainer{URLService: URLService}
	handler := handlerContainer.GetFullURLHandler()

	tests := []struct {
		Name                 string
		Path                 string
		ExpectedResponseCode int
	}{
		{
			"absent URL id",
			"http://localhost/absent-id",
			http.StatusNotFound,
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

type expectedAPIResponse struct {
	Result string `json:"result"`
}

func TestCreateShortURLApiHandler(t *testing.T) {
	db, err := storage.NewStorage("")
	require.NoError(t, err)
	URLService := service.Service{Storage: db}
	handlerContainer := HandlerContainer{URLService: URLService}
	handler := handlerContainer.CreateShortURLApiHandler()

	tests := []struct {
		Name                 string
		Method               string
		Path                 string
		Body                 string
		ExpectedResponseCode int
	}{
		{
			"POST with invalid body",
			http.MethodPost,
			"http://localhost/",
			"invalid body",
			http.StatusBadRequest,
		},
		{
			"POST with invalid URL",
			http.MethodPost,
			"http://localhost/",
			`{"url": "invalid URL"}`,
			http.StatusBadRequest,
		},
		{
			"POST happy pass",
			http.MethodPost,
			"http://localhost/",
			`{"url": "http://some_test/url/test"}`,
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
			if tc.ExpectedResponseCode == http.StatusCreated {
				contentHeader := result.Header.Get("content-type")
				expectedContentType := "application/json; charset=utf-8"
				assert.Equal(t, expectedContentType, contentHeader)
				var responsePayload expectedAPIResponse
				responseBody, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				err = json.Unmarshal(responseBody, &responsePayload)
				require.NoError(t, err)
				_, err = url.ParseRequestURI(responsePayload.Result)
				require.NoError(t, err)
			}
		})
	}
}
