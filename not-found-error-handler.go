package gfs

import (
	"html/template"
	"log"
	"net/http"
)

const (
	//language=html
	NotFoundHtml string = `<!DOCTYPE html>
<html>
<head>
</head>
<body>
<h1>Not found</h1>
<p>Unable to find requested resource: {{.Path}}</p>
</body>
</html>`
)

type notFound struct {
	Path string `json:"path" xml:"path"`
}

type NotFoundHandler struct {
	responseHandler
	htmlTemplate *template.Template
}

func GetNotFoundHandler() (h *NotFoundHandler, err error) {
	t := template.New("Not Found Html Template")
	t, err = t.Parse(NotFoundHtml)
	if err != nil {
		return nil, err
	}
	h = &NotFoundHandler{
		htmlTemplate: t,
	}

	return h, nil
}

// Write a not found message to the response.
func (h *NotFoundHandler) Handle(writer http.ResponseWriter, p, format string) {
	response := notFound{p}

	err := h.responseHandler.WriteResponse(writer, http.StatusNotFound, h.htmlTemplate, format, response)
	if err != nil {
		log.Println("Something went wrong when responding", err.Error())
	}
}
