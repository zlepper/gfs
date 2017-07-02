package gfs

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	//language=html
	FileResponseHtml string = `<!DOCTYPE html>
<html>
<head>
</head>
<body>
<h1><a href="{{.Path}}" download>{{.Name}}</a></h1>
<hr/>
<table>
    <tbody>
    <tr>
        <th>Name:</th>
        <td>{{.Name}}</td>
    </tr>
    <tr>
        <th>Size:</th>
        <td>{{.Size}}</td>
    </tr>
    <tr>
        <th>Last modification:</th>
        <td>{{.LastModificationTime}}</td>
    </tr>
    </tbody>
</table>
<hr/>
</body>
</html>`
)

type FileStats struct {
	// The name of the file
	Name string `json:"name" xml:"name"`
	// The path to the file. Relative to the serve root
	Path string `json:"path" xml:"path"`
	// The size of the file. 0 if a directory
	Size int64 `json:"size,omitempty" xml:"size,omitempty"`
	// The last time this file was modified
	LastModificationTime time.Time `json:"last_modification_time" xml:"last_modification_time"`
}

type FileResponseHandler struct {
	responseHandler
	htmlTemplate *template.Template
}

func GetFileResponseHandler() (h *FileResponseHandler, err error) {
	t := template.New("File Response Html Template")
	t, err = t.Parse(FileResponseHtml)
	if err != nil {
		return nil, err
	}
	h = &FileResponseHandler{
		htmlTemplate: t,
	}

	return h, nil
}

// Write a not found message to the response.
func (h *FileResponseHandler) Handle(writer http.ResponseWriter, fullpath, p string, format string) error {
	if format == "" {
		file, err := os.Open(fullpath)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	}

	stats, err := GetFileStats(fullpath, p)
	if err != nil {
		return err
	}

	err = h.responseHandler.WriteResponse(writer, http.StatusOK, h.htmlTemplate, format, stats)
	if err != nil {
		log.Println("Something went wrong when responding", err.Error())
	}
	return nil
}

// Gets the stats about a specific file
func GetFileStats(fullpath, p string) (*FileStats, error) {
	stats, err := os.Stat(fullpath)
	if err != nil {
		return nil, err
	}

	return &FileStats{
		Name:                 stats.Name(),
		Path:                 p,
		Size:                 stats.Size(),
		LastModificationTime: stats.ModTime(),
	}, nil

}
