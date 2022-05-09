package storage

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStorage(t *testing.T) {
	type args struct {
		storageFile string
	}
	tests := []struct {
		name string
		args args
		want Repository
	}{
		{"memory", args{""}, MemoryStorage(make(map[string]string))},
		{"file", args{"test.txt"}, FileStorage{"test.txt"}},
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
		{"existed storage file", fields{"test.txt"}, args{"123", "test_url"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := FileStorage{filename: tt.fields.filename}
			if tt.newFile == false {
				_, err := os.Create(tt.fields.filename)
				require.NoError(t, err)
			}
			err := s.WriteURL(tt.args.urlID, tt.args.fullURL)
			require.NoError(t, err)
			storageFile, err := os.Open(tt.fields.filename)
			require.NoError(t, err)
			decoder := json.NewDecoder(storageFile)
			savedURL := &URLInfo{}
			err = decoder.Decode(&savedURL)
			require.NoError(t, err)
			require.Equal(t, URLInfo{tt.args.urlID, tt.args.fullURL}, *savedURL)
			os.Remove(tt.fields.filename)
		})
	}
}
