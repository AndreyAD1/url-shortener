package server_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndreyAD1/url-shortener/internal/app/config"
	"github.com/AndreyAD1/url-shortener/internal/app/server"
)

const testURL = "https://github.com/AndreyAD1"

var testConfig = config.StartupConfig{
	ServerAddress:  "localhost:8080",
	BaseURL:        "http://localhost:8080",
	ShortURLLength: 10,
}

func getTestServer(t *testing.T) *httptest.Server {
	listener, err := net.Listen("tcp", testConfig.ServerAddress)
	require.NoError(t, err)
	server := &httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: server.GetHandler(testConfig)},
	}
	server.Start()
	return server
}

func Test_GetShortURL(t *testing.T) {
	server := getTestServer(t)
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	requestBody := bytes.NewBufferString(testURL)
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		server.URL,
		requestBody,
	)
	require.NoError(t, err)
	client := &http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()
	require.Equal(t, http.StatusCreated, response.StatusCode)
	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	returnedURL, err := url.ParseRequestURI(string(body))
	require.NoError(t, err)
	assert.Equal(t, "http", returnedURL.Scheme)
	assert.Equal(t, testConfig.ServerAddress, returnedURL.Host)
	assert.Equal(t, testConfig.ShortURLLength, len(returnedURL.Path[1:]))
}

func getShortURL(t *testing.T, server *httptest.Server) string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	requestBody := bytes.NewBufferString(testURL)
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		server.URL,
		requestBody,
	)
	require.NoError(t, err)
	client := &http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()
	require.Equal(t, http.StatusCreated, response.StatusCode)
	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	return string(body)
}

func Test_GetFullURL(t *testing.T) {
	server := getTestServer(t)
	defer server.Close()
	shortURL := getShortURL(t, server)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		shortURL,
		nil,
	)
	require.NoError(t, err)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	response, err := client.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()
	require.Equal(t, http.StatusTemporaryRedirect, response.StatusCode)
	require.Equal(t, testURL, response.Header.Get("Location"))
}

func Test_GetFullURL_error(t *testing.T) {
	server := getTestServer(t)
	defer server.Close()
	tests := []struct {
		Name                 string
		Method               string
		Path                 string
		ExpectedResponseCode int
	}{
		{
			"GET without id",
			http.MethodGet,
			server.URL,
			http.StatusMethodNotAllowed,
		},
		{
			"POST with id",
			http.MethodPost,
			server.URL + "/123",
			http.StatusMethodNotAllowed,
		},
		{
			"GET with invalid path",
			http.MethodGet,
			server.URL + "/test/path",
			http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(
				context.Background(),
				5*time.Second,
			)
			defer cancel()
			request, err := http.NewRequestWithContext(
				ctx,
				tc.Method,
				tc.Path,
				nil,
			)
			require.NoError(t, err)
			client := &http.Client{}
			response, err := client.Do(request)
			require.NoError(t, err)
			defer response.Body.Close()
			require.Equal(t, tc.ExpectedResponseCode, response.StatusCode)
		})
	}
}
