package gfs

import (
	"html/template"
	"log"
	"net/http"
)

const (
	//language=html
	InternalServerErrorHtml string = `<!DOCTYPE html>
<html>
<head>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Something went wrong when processing the request:<p>
<pre>
    {{.Error}}
</pre>
</body>
</html>`
)

type internalServerError struct {
	Error string `json:"error" xml:"error"`
}

type InternalServerErrorHandler struct {
	responseHandler
	htmlTemplate *template.Template
}

func GetInternalServerErrorHandler() (h *InternalServerErrorHandler, err error) {
	t := template.New("Internal Server Error Template")
	t, err = t.Parse(InternalServerErrorHtml)
	if err != nil {
		return nil, err
	}
	h = &InternalServerErrorHandler{
		htmlTemplate: t,
	}

	return h, nil
}

// Write a not found message to the response.
func (h *InternalServerErrorHandler) Handle(writer http.ResponseWriter, err error, format string) {
	response := internalServerError{err.Error()}

	err = h.responseHandler.WriteResponse(writer, http.StatusInternalServerError, h.htmlTemplate, format, response)
	if err != nil {
		log.Println("Something went wrong when responding", err.Error())
	}
}
