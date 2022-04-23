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

func getTestServer(t *testing.T) *httptest.Server {
	listener, err := net.Listen("tcp", config.ServerAddress)
	require.NoError(t, err)
	server := &httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: server.GetHandler()},
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
	assert.Equal(t, config.ServerAddress, returnedURL.Host)
	assert.Equal(t, config.ShortURLLength, len(returnedURL.Path[1:]))
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
