package storage

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStorage(t *testing.T) {
	URLItem := URLInfo{"URL ID", "test URL"}

	type args struct {
		storageFile string
	}
	tests := []struct {
		name    string
		args    args
		newFile bool
		want    Repository
	}{
		{"memory", args{""}, true, &MemoryStorage{storage: make(map[string]string)}},
		{
			"file",
			args{"test.txt"},
			true,
			&FileStorage{filename: "test.txt", storage: make(map[string]string)},
		},
		{
			"file",
			args{"test.txt"},
			false,
			&FileStorage{
				filename: "test.txt",
				storage:  map[string]string{URLItem.ID: URLItem.URL},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.newFile == false {
				file, err := os.Create(tt.args.storageFile)
				require.NoError(t, err)
				err = json.NewEncoder(file).Encode(URLItem)
				require.NoError(t, err)
			}
			defer os.Remove(tt.args.storageFile)
			got, err := NewStorage(tt.args.storageFile)
			require.NoError(t, err)
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
		name   string
		fields fields
		args   args
	}{
		{"new storage file", fields{"test.txt"}, args{"123", "test_url"}},
		{"existing storage file", fields{"test.txt"}, args{"123", "test_url"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewStorage(tt.fields.filename)
			require.NoError(t, err)
			defer os.Remove(tt.fields.filename)
			err = s.WriteURL(tt.args.urlID, tt.args.fullURL)
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
			if tt.newFile == false {
				file, err := os.Create(tt.fields.filename)
				require.NoError(t, err)
				for _, URLItem := range tt.savedURLs {
					err = json.NewEncoder(file).Encode(URLItem)
					require.NoError(t, err)
				}
			}
			defer os.Remove(tt.fields.filename)
			s, err := NewStorage(tt.fields.filename)
			require.NoError(t, err)

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
