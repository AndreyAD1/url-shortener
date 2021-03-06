package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AndreyAD1/url-shortener/internal/app/handlers"
)

type expectedResponse struct {
	Result string `json:"result"`
}

func Test_GetShortURLviaAPI(t *testing.T) {
	server := getTestServer(t)
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	payload := handlers.CreateShortURLRequest{URL: testURL}
	requestBytes, err := json.Marshal(payload)
	require.NoError(t, err)
	requestBody := bytes.NewBuffer(requestBytes)
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		server.URL+"/api/shorten",
		requestBody,
	)
	require.NoError(t, err)
	client := &http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()
	require.Equal(t, http.StatusCreated, response.StatusCode)
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	var responsePayload expectedResponse
	err = json.Unmarshal(body, &responsePayload)
	require.NoError(t, err)
	returnedURL, err := url.ParseRequestURI(responsePayload.Result)
	require.NoError(t, err)
	expectedURL, err := url.ParseRequestURI(testConfig.BaseURL)
	require.NoError(t, err)
	assert.Equal(t, expectedURL.Scheme, returnedURL.Scheme)
	assert.Equal(t, expectedURL.Host, returnedURL.Host)
	assert.Equal(t, testConfig.ShortURLLength, len(returnedURL.Path[1:]))
}
