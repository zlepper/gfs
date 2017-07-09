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
	config             *Config
	clientErrorHandler *ClientErrorHandler
}

var (
	ErrNoUploadingUp      error = errors.New("Unable to upload up outside the <serve> directory.")
	ErrNoFilenameProvided       = errors.New("No filename provided. Cannot accept upload")
	ErrUnknownContentType       = errors.New("Unknown upload content type. Cannot proceed.")
)

func isUploadClientError(err error) bool {
	return err == ErrNoUploadingUp ||
		err == ErrNoFilenameProvided ||
		err == ErrUnknownContentType
}

func (h *UploadHandler) Handle(writer http.ResponseWriter, request *http.Request, responseFormat string) error {

	ct := getContentType(request)
	if ct == "multipart/form-data" {

		uploadPath := request.FormValue("path")

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

				return h.uploadFile(path.Join(uploadPath, fileRef.Filename), file)
			}(i)
			if err != nil {
				return err
			}
		}

		http.Redirect(writer, request, uploadPath, http.StatusFound)

		return nil
	} else if ct == FormatOctetStream {
		filename := request.URL.Query().Get("filename")
		if filename == "" {
			return ErrNoFilenameProvided
		}

		defer request.Body.Close()

		return h.uploadFile(filename, request.Body)

	} else {
		return ErrUnknownContentType
	}
}

func (h *UploadHandler) uploadFile(filename string, file io.Reader) error {
	outputPath := path.Join(h.config.Serve, filename)
	log.Println("outputPath", outputPath)

	// Ensure that it's not possible to upload "upwards" in the tree
	if !strings.HasPrefix(outputPath, h.config.Serve) {
		return ErrNoUploadingUp
	}

	err := os.MkdirAll(path.Dir(outputPath), os.ModePerm)
	if err != nil {
		return err
	}

	dst, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	return err
}

func GetUploadHandler(config *Config) (*UploadHandler, error) {
	chl, err := GetClientErrorHandler()
	if err != nil {
		return nil, err
	}
	return &UploadHandler{
		config:             config,
		clientErrorHandler: chl,
	}, nil
}
