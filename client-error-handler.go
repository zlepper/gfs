package gfs

import (
	"html/template"
	"log"
	"net/http"
)

const (
	//language=html
	ClientErrorHtml string = `<!DOCTYPE html>
<html>
<head>
</head>
<body>
<h1>Invalid request</h1>
<p>Invalid request:</p>
<pre>
{{.Error}}
</pre>
</body>
</html>`
)

type invalidRequest struct {
	Error string `json:"error" xml:"error"`
}

type ClientErrorHandler struct {
	responseHandler
	htmlTemplate *template.Template
}

func GetClientErrorHandler() (h *ClientErrorHandler, err error) {
	t := template.New("Client Error Html Template")
	t, err = t.Parse(ClientErrorHtml)
	if err != nil {
		return nil, err
	}
	h = &ClientErrorHandler{
		htmlTemplate: t,
	}

	return h, nil
}

// Write a not found message to the response.
func (h *ClientErrorHandler) Handle(writer http.ResponseWriter, err error, format string, errorCode int) {
	response := invalidRequest{err.Error()}

	err = h.responseHandler.WriteResponse(writer, errorCode, h.htmlTemplate, format, response)
	if err != nil {
		log.Println("Something went wrong when responding", err.Error())
	}
}
