package gfs

import (
	"log"
	"net/http"
	"strings"
)

func getResponseFormat(request *http.Request) string {
	accepts := request.Header.Get("accept")
	if accepts == "" {
		return ""
	}

	acceptOptions := strings.Split(accepts, ",")

	for _, option := range acceptOptions {
		switch option {
		case FormatHtml:
			return FormatHtml
		case "text/json":
			fallthrough
		case FormatJson:
			return FormatJson
		case "text/xml":
			fallthrough
		case FormatXml:
			return FormatXml
		}
	}

	log.Println("Could not find known format in accepts. Hoping for the best. Request accept header:", accepts)

	return ""
}

func getContentType(request *http.Request) string {
	return request.Header.Get("Content-Type")
}