package gfs

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type UploadHandler struct {
	responseHandler
	config *Config
}

var (
	ErrNoUploadingUp error = errors.New("Unable to upload up outside the <serve> directory.")
)

func (h *UploadHandler) Handle(writer http.ResponseWriter, request *http.Request) error {
	uploadPath := request.FormValue("path")
	outputPath := path.Join(h.config.Serve, uploadPath)
	log.Println("outputPath", outputPath)

	// Ensure that it's not possible to upload "upwards" in the tree
	if !strings.HasPrefix(outputPath, h.config.Serve) {
		return ErrNoUploadingUp
	}

	err := request.ParseMultipartForm(1 << 20) // 1MB
	if err != nil {
		return err
	}

	files, ok := request.MultipartForm.File["uploadfiles"]
	if !ok {
		log.Println("No upload files")
		log.Println(request.MultipartForm)
	}
	log.Println(files)
	for i := range files {
		err := func(i int) error {
			fileRef := files[i]
			log.Println("Handling file", fileRef.Filename)
			file, err := fileRef.Open()
			if err != nil {
				return err
			}
			defer file.Close()

			dst, err := os.Create(path.Join(outputPath, fileRef.Filename))
			if err != nil {
				return err
			}
			defer dst.Close()

			_, err = io.Copy(dst, file)
			return err
		}(i)
		if err != nil {
			return err
		}
	}

	http.Redirect(writer, request, uploadPath, http.StatusFound)

	return nil
}

func GetUploadHandler(config *Config) (*UploadHandler, error) {
	return &UploadHandler{
		config: config,
	}, nil
}
