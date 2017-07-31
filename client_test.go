package gfs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_Login(t *testing.T) {
	a := assert.New(t)

	c, err := NewClient("http://localhost:28080", "admin", "password")
	if a.NoError(err) {
		a.NotNil(c)
	}
}

func TestClient_UploadFiles(t *testing.T) {
	a := assert.New(t)

	c, err := NewClient("http://localhost:28080", "admin", "password")
	if a.NoError(err) {
		a.NotNil(c)

		f, err := NewUploadFileFromDisk("client.go", "test-path")

		if a.NoError(err) {
			a.Equal("client.go", f.Filename)

			err = c.UploadFile(f)

			a.NoError(err)
		}
	}
}
