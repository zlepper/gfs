package gfs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_Login(t *testing.T) {
	a := assert.New(t)

	c, err := NewClient("http://localhost:8080", "admin", "password")
	if a.NoError(err) {
		a.NotNil(c)
	}
}

func TestClient_UploadFiles(t *testing.T) {
	a := assert.New(t)

	c, err := NewClient("http://localhost:8080", "admin", "password")
	if a.NoError(err) {
		a.NotNil(c)

		f, err := NewUploadFileFromDisk("client.go")

		if a.NoError(err) {
			err = c.UploadFile(f)

			a.NoError(err)
		}
	}
}
