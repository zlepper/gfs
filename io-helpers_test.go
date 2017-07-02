package gfs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDirectoryStats(t *testing.T) {
	a := assert.New(t)
	stats, err := GetDirectoryStats("./handlers", ".")
	if a.NoError(err) {
		a.Equal("/handlers", stats.Path)
		a.Equal("handlers", stats.Name)
		if a.NotNil(stats.Entries) {
			entries := stats.Entries
			a.Equal("client-error-handler.go", entries[0].Name)
			a.Equal("response-handler.go", entries[len(entries)-1].Name)
		}
	}
}
