package gfs

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

const (
	//language=html
	LoginHtml = `
<h2>Login</h2>
<form action="/login" method="post">
    <label for="usernameInput">Username</label>
    <input name="username" id="usernameInput" type="text" required />
    <label for="passwordInput">Password</label>
    <input name="password" id="passwordInput" type="password" required />
    <input type="hidden" name="redirectTo" value="{{.Path}}" />
    <button type="submit">Login</button>
</form>`

	//language=html
	UploadHtml string = `<form enctype="multipart/form-data" name="uploadFilesForm" id="uploadFilesForm" action="/upload" method="post">
    <input type="file" multiple="multiple" name="uploadfiles"/>
    <input type="hidden" name="path" value="{{.Path}}"/>
    <button type="submit" form="uploadFilesForm">Upload</button>
</form>
	`
	//language=html
	DirectoryResponseHtml string = `<!DOCTYPE html>
<html>
<head>
<title>{{.Path}}</title>
</head>
<body>
<h1><a href="{{.Path}}">{{.Name}}</a> <small>last modified: {{.LastModificationTime}}</small></h1>
<hr />
<table>
    <thead>
        <tr>
            <th>Name</td>
            <th>Size</td>
            <th>Last modified</td>
        </tr>
    </thead>
    <tbody>
		{{range .Entries}}
		<tr>
			<td><a href="{{.Path}}">{{.Name}}</a></td>
			<td>
			{{if .IsDirectory}}
				&lt;Dir&gt;
			{{else}}
				{{.Size}}
			{{end}}
			</td>
			<td>
			{{.LastModificationTime}}
			</td>
		</tr>
		{{else}}
		<tr>
			<td colspan="3">No entries</td>
		</tr>
		{{end}}
    </tbody>
</table>
<hr />
{{if .Authorized}}
	{{template "upload" .}}
{{else}}
	{{template "login" .}}
{{end}}
</body>
</html>`
)

type DirectoryEntry struct {
	// The name of the file
	Name string `json:"name" xml:"name"`
	// The path to the file. Relative to the serve root
	Path string `json:"path" xml:"path"`
	// The size of the file. 0 if a directory
	Size int64 `json:"size,omitempty" xml:"size,omitempty"`
	// True is directory. Probably false if not
	IsDirectory bool `json:"is_directory" xml:"is_directory"`
	// The last time this file was modified
	LastModificationTime time.Time `json:"last_modification_time" xml:"last_modification_time"`
}

// Simple stats about a directory
type DirectoryStats struct {
	// The name of the directory
	Name string `json:"name" xml:"name"`
	// The path to the directory. Relative to the serve root
	Path string `json:"path" xml:"path"`
	// The last time this file was modified
	LastModificationTime time.Time `json:"last_modification_time" xml:"last_modification_time"`
	// The content directory available in this directory
	Entries []DirectoryEntry `json:"entries" xml:"entries"`
	// Indicates if the request is authorized
	Authorized bool `json:"authorized" xml:"authorized"`
}

type DirectoryResponseHandler struct {
	responseHandler
	htmlTemplate   *template.Template
	loginTemplate  *template.Template
	uploadTemplate *template.Template
}

func GetDirectoryResponseHandler() (h *DirectoryResponseHandler, err error) {

	t := template.New("Directory Response Html Template")

	loginTemplate := t.New("login")
	if loginTemplate, err = loginTemplate.Parse(LoginHtml); err != nil {
		return nil, err
	}

	uploadTemplate := t.New("upload")
	if uploadTemplate, err = uploadTemplate.Parse(UploadHtml); err != nil {
		return nil, err
	}

	t, err = t.Parse(DirectoryResponseHtml)
	if err != nil {
		return nil, err
	}
	h = &DirectoryResponseHandler{
		htmlTemplate:   t,
		uploadTemplate: uploadTemplate,
		loginTemplate:  loginTemplate,
	}

	return h, nil
}

// Write a not found message to the response.
func (h *DirectoryResponseHandler) Handle(writer http.ResponseWriter, stats *DirectoryStats, format string) {
	err := h.responseHandler.WriteResponse(writer, http.StatusOK, h.htmlTemplate, format, stats)
	if err != nil {
		log.Println("Something went wrong when responding", err.Error())
	}
}
