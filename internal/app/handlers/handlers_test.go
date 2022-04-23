package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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
		name string
		Method string
		Path string
		Body string
		expectedResponseCode int
		expectedResponseBody string
	}{
		{
			"invalid request method",
			"PUT",
			"whatever",
			"",
			http.StatusMethodNotAllowed,
			"Only GET and POST methods are allowed",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
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

			assert.Equal(t, tc.expectedResponseCode, result.StatusCode)
			body, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			bodyString := strings.TrimSpace(string(body))
			assert.Equal(t, tc.expectedResponseBody, bodyString)
		})
	}
}
