package server_test

import (
	"bytes"
	"compress/gzip"
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

func Test_SendGzipEncodedRequest(t *testing.T) {
	server := getTestServer(t)
	defer server.Close()
	testCases := []struct {
		name string
		requestIsCompressed bool
		expectCompressedResponse bool
	} {
		{
			"Compressed request, full response",
			true,
			false,
		},
		{
			"Full request, compressed response",
			false,
			true,
		},
		{
			"Compressed request, compressed response",
			true,
			true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(
				context.Background(), 
				1*time.Second,
			)
			defer cancel()
			payload := handlers.CreateShortURLRequest{URL: testURL}
			requestBytes, err := json.Marshal(payload)
			require.NoError(t, err)
			
			var requestBodyReader bytes.Buffer
			if testCase.requestIsCompressed {
				gzipWriter := gzip.NewWriter(&requestBodyReader)
				_, err = gzipWriter.Write(requestBytes)
				require.NoError(t, err)
				gzipWriter.Close()
			} else {
				_, err = requestBodyReader.Write(requestBytes)
				require.NoError(t, err)
			}

			
			request, err := http.NewRequestWithContext(
				ctx,
				http.MethodPost,
				server.URL+"/api/shorten",
				&requestBodyReader,
			)
			require.NoError(t, err)
			if testCase.requestIsCompressed {
				request.Header.Set("Content-Encoding", "gzip")
			}
			if testCase.expectCompressedResponse {
				request.Header.Set("Accept-Encoding", "gzip")
			}
			client := &http.Client{}
			response, err := client.Do(request)
			require.NoError(t, err)
			defer response.Body.Close()
			
			require.Equal(t, http.StatusCreated, response.StatusCode)
			
			var body []byte
			if testCase.expectCompressedResponse {
				require.Equal(t, "gzip", response.Header.Get("Content-Encoding"))
				gz, err := gzip.NewReader(response.Body)
				require.NoError(t, err)
				defer gz.Close()
				body, err = io.ReadAll(gz)
				require.NoError(t, err)
			} else {
				body, err = io.ReadAll(response.Body)
				require.NoError(t, err)
			}

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
		})
	}
}
