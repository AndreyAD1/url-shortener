package storage

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStorage(t *testing.T) {
	_, err := os.Create("test.txt")
	require.NoError(t, err)
	defer os.Remove("test.txt")

	type args struct {
		storageFile string
	}
	tests := []struct {
		name string
		args args
		want Repository
	}{
		{"memory", args{""}, &MemoryStorage{storage: make(map[string]string)}},
		{
			"file",
			args{"test.txt"},
			&FileStorage{filename: "test.txt", storage: make(map[string]string)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewStorage(tt.args.storageFile)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFileStorage_WriteURL(t *testing.T) {
	type fields struct {
		filename string
	}
	type args struct {
		urlID   string
		fullURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		newFile bool
	}{
		{"new storage file", fields{"test.txt"}, args{"123", "test_url"}, false},
		{"existing storage file", fields{"test.txt"}, args{"123", "test_url"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := FileStorage{filename: tt.fields.filename}
			if tt.newFile == false {
				_, err := os.Create(tt.fields.filename)
				require.NoError(t, err)
			}
			defer os.Remove(tt.fields.filename)
			err := s.WriteURL(tt.args.urlID, tt.args.fullURL)
			require.NoError(t, err)
			storageFile, err := os.Open(tt.fields.filename)
			require.NoError(t, err)
			decoder := json.NewDecoder(storageFile)
			savedURL := &URLInfo{}
			err = decoder.Decode(&savedURL)
			require.NoError(t, err)
			require.Equal(t, URLInfo{tt.args.urlID, tt.args.fullURL}, *savedURL)
		})
	}
}

func TestFileStorage_GetURL(t *testing.T) {
	type fields struct {
		filename string
	}
	type args struct {
		urlID string
	}
	tests := []struct {
		name        string
		newFile     bool
		fields      fields
		args        args
		savedURLs   []URLInfo
		URLFound    bool
		expectedURL string
	}{
		{"no storage file", true, fields{"test.txt"}, args{"123"}, []URLInfo{}, false, ""},
		{"no saved URLs", false, fields{"test.txt"}, args{"123"}, []URLInfo{}, false, ""},
		{
			"one saved URL - not found",
			false,
			fields{"test.txt"},
			args{"123"},
			[]URLInfo{{"999", "some test url"}},
			false,
			"",
		},
		{
			"two saved URLs - not found",
			false,
			fields{"test.txt"},
			args{"123"},
			[]URLInfo{{"999", "some test url"}, {"111", "another test url"}},
			false,
			"",
		},
		{
			"two saved URLs - found first",
			false,
			fields{"test.txt"},
			args{"999"},
			[]URLInfo{{"999", "some test url"}, {"111", "another test url"}},
			true,
			"some test url",
		},
		{
			"two saved URLs - found second",
			false,
			fields{"test.txt"},
			args{"111"},
			[]URLInfo{{"999", "some test url"}, {"111", "another test url"}},
			true,
			"another test url",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := FileStorage{filename: tt.fields.filename}
			if tt.newFile == false {
				_, err := os.Create(tt.fields.filename)
				require.NoError(t, err)
			}
			defer os.Remove(tt.fields.filename)
			for _, savedURL := range tt.savedURLs {
				err := s.WriteURL(savedURL.ID, savedURL.URL)
				require.NoError(t, err)
			}
			returnedURL, err := s.GetURL(tt.args.urlID)
			require.NoError(t, err)
			if tt.URLFound {
				require.Equal(t, tt.expectedURL, *returnedURL)
			} else {
				require.Nil(t, returnedURL)
			}
		})
	}
}
