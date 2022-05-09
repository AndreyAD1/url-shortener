package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStorage(t *testing.T) {
	type args struct {
		storageFile string
	}
	tests := []struct {
		name    string
		args    args
		want    Repository
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
