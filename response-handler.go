package gfs

import (
	"encoding/json"
	"encoding/xml"
	"html/template"
	"log"
	"net/http"
)

// A basic handler that suplies basic logic for writing responses
type responseHandler struct {
}

// Takes care of writing a response in the correct format
func (r *responseHandler) WriteResponse(writer http.ResponseWriter, statusCode int, htmlTemplate *template.Template, format string, response interface{}) error {

	switch format {
	default:
		log.Println("Unknown response format", format)
		fallthrough
	case FormatHtml:
		writer.Header().Set("content-type", FormatHtml)
		writer.WriteHeader(statusCode)
		return htmlTemplate.Execute(writer, response)
	case FormatJson:
		writer.Header().Set("content-type", FormatJson)
		writer.WriteHeader(statusCode)
		return json.NewEncoder(writer).Encode(response)
	case FormatXml:
		writer.Header().Set("content-type", FormatXml)
		writer.WriteHeader(statusCode)
		return xml.NewEncoder(writer).Encode(response)
	}
}
