package gfs

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

const (
	//language=html
	LoginFailedHtml string = `
	<p style="color: red">Login failed: {{.Error}}</p>
	{{template "login" .}}
	`
)

type LoginRequest struct {
	Username string `json:"username" xml:"username"`
	Password string `json:"password" xml:"password"`
}

type AuthoizationFailedResponse struct {
	Path  string `json:"redirect_path" xml:"redirect_path"`
	Error string `json:"error" xml:"error"`
}

type AuthorizationSuccessResponse struct {
	Token string `json:"token" xml:"token"`
}

type AuthorizationResponse struct {
	Path  string `json:"redirect_path" xml:"redirect_path"`
	Error string `json:"error" xml:"error"`
	Token string `json:"token" xml:"token"`
}

type AuthorizationHandler struct {
	responseHandler
	config              *Config
	loginFailedTemplate *template.Template
}

func (h *AuthorizationHandler) Login(writer http.ResponseWriter, request *http.Request, format string) error {

	var username, password, redirectPath string
	contentType := getContentType(request)
	switch contentType {
	case FormatXFormUrlEncoded:
		username = request.FormValue("username")
		password = request.FormValue("password")
		redirectPath = request.FormValue("redirectTo")
	case FormatJson:
		var loginRequest LoginRequest
		err := json.NewDecoder(request.Body).Decode(&loginRequest)
		if err != nil {
			return err
		}
		username = loginRequest.Username
		password = loginRequest.Password
	case FormatXml:
		var loginRequest LoginRequest
		err := xml.NewDecoder(request.Body).Decode(&loginRequest)
		if err != nil {
			return err
		}
		username = loginRequest.Username
		password = loginRequest.Password
	default:
		return errors.New(fmt.Sprintf("Unknown request format '%s'. Accepted types are: '%s', '%s' and '%s'", contentType, FormatXFormUrlEncoded, FormatJson, FormatXml))
	}

	if h.config.Username == username {
		matches, err := CheckPassword(password, h.config.Password)
		if err != nil {
			fail := AuthoizationFailedResponse{Path: redirectPath, Error: err.Error()}
			return h.responseHandler.WriteResponse(writer, http.StatusInternalServerError, h.loginFailedTemplate, format, fail)
		}
		if !matches {
			fail := AuthoizationFailedResponse{Path: redirectPath, Error: "Invalid username or password"}
			return h.responseHandler.WriteResponse(writer, http.StatusBadRequest, h.loginFailedTemplate, format, fail)
		}

		token, err := GetToken([]byte(h.config.Secret))
		if err != nil {
			fail := AuthoizationFailedResponse{Path: redirectPath, Error: err.Error()}
			return h.responseHandler.WriteResponse(writer, http.StatusInternalServerError, h.loginFailedTemplate, format, fail)
		}

		if format == FormatXml || format == FormatJson {
			response := AuthorizationSuccessResponse{Token: token}
			return h.WriteResponse(writer, http.StatusOK, nil, format, response)
		} else {
			cookie := &http.Cookie{
				Name:    "token",
				Value:   token,
				Path:    "/",
				Expires: time.Now().Add(31 * 24 * time.Hour),
				MaxAge:  31 * 24 * 60 * 60,
			}

			http.SetCookie(writer, cookie)
			http.Redirect(writer, request, redirectPath, http.StatusFound)
		}
		return nil
	}
	fail := AuthoizationFailedResponse{Path: redirectPath, Error: "Invalid username or password"}
	return h.responseHandler.WriteResponse(writer, http.StatusBadRequest, h.loginFailedTemplate, format, fail)
}

// Checks if the request is authenticated. Returns nil if request is authenticated
func (h *AuthorizationHandler) CheckAuthenticated(request *http.Request) error {
	token := request.Header.Get("gfs-token")
	if token == "" {
		cookie, err := request.Cookie("token")
		if err != nil {
			return err
		}

		token = cookie.Value
	}

	err := GetTokenData(token, []byte(h.config.Secret), &struct{}{})
	return err
}

func GetAuthorizationHandler(config *Config) (*AuthorizationHandler, error) {
	loginFailedTemplate := template.New("loginFailed")
	var err error

	loginTemplate := loginFailedTemplate.New("login")
	loginTemplate, err = loginTemplate.Parse(LoginHtml)
	if err != nil {
		return nil, err
	}

	loginFailedTemplate, err = loginFailedTemplate.Parse(LoginFailedHtml)
	if err != nil {
		return nil, err
	}

	h := &AuthorizationHandler{
		config:              config,
		loginFailedTemplate: loginFailedTemplate,
	}

	return h, nil
}
