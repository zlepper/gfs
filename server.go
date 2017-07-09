package gfs

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

const (
	// The running version of GFS
	GFSVersion string = "0.0.2"
)

func RunServer(config *Config) {
	go checkForUpdates() // Check for updates on startup

	handlerFunc, err := getHandler(config)
	if err != nil {
		log.Fatalln(err)
	}
	http.HandleFunc("/", handlerFunc)

	loginHandlerFunc, err := getLoginHandler(config, handlerFunc)
	if err != nil {
		log.Fatalln(err)
	}
	http.HandleFunc("/login", loginHandlerFunc)

	uploadHandlerFunc, err := getUploadHandlerFunc(config, handlerFunc)
	if err != nil {
		log.Fatalln(err)
	}
	http.HandleFunc("/upload", uploadHandlerFunc)

	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}

func getHandler(config *Config) (f http.HandlerFunc, err error) {

	notFoundHandler, err := GetNotFoundHandler()
	if err != nil {
		return nil, err
	}
	internalServerErrorHandler, err := GetInternalServerErrorHandler()
	if err != nil {
		return nil, err
	}
	clientErrorHandler, err := GetClientErrorHandler()
	if err != nil {
		return nil, err
	}
	directoryResponseHandler, err := GetDirectoryResponseHandler()
	if err != nil {
		return nil, err
	}
	fileResponserHandler, err := GetFileResponseHandler()
	if err != nil {
		return nil, err
	}
	authorizationHandler, err := GetAuthorizationHandler(config)
	if err != nil {
		return nil, err
	}

	f = func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("gfs-version", GFSVersion)
		responseFormat := getResponseFormat(request)

		if request.Method == "GET" {
			authorized := false
			err := authorizationHandler.CheckAuthenticated(request)
			if err == nil {
				authorized = true
			}
			if config.LoginRequiredForRead && !authorized {
				log.Println("authentication error", err)
				clientErrorHandler.Handle(writer, errors.New("Not authenticated"), responseFormat, http.StatusUnauthorized)
				return
			}

			p := request.URL.Path
			fullpath := path.Join(config.Serve, p)

			directory, err := isDirectory(fullpath)
			if err != nil {
				if os.IsNotExist(err) {
					notFoundHandler.Handle(writer, p, responseFormat)
					return
				}
				log.Println("Error when detecting directory", err)
				internalServerErrorHandler.Handle(writer, err, responseFormat)
				return
			}

			if directory {
				stats, err := GetDirectoryStats(fullpath, p)
				if err != nil {
					internalServerErrorHandler.Handle(writer, err, responseFormat)
				}
				stats.Authorized = authorized
				directoryResponseHandler.Handle(writer, stats, responseFormat)
			} else {
				err := fileResponserHandler.Handle(writer, fullpath, p, responseFormat)
				if err != nil {
					log.Println("Something went wrong when serving file:", err.Error())
				}
			}
			return
		}

		clientErrorHandler.Handle(writer, errors.New(fmt.Sprintf("Unsupported method: '%s'", request.Method)), responseFormat, http.StatusMethodNotAllowed)
	}
	return f, nil
}

func getLoginHandler(config *Config, defaultHandler http.HandlerFunc) (http.HandlerFunc, error) {
	authorizationHandler, err := GetAuthorizationHandler(config)
	if err != nil {
		return nil, err
	}

	f := func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("gfs-version", GFSVersion)

		if request.Method == "POST" {

			responseFormat := getResponseFormat(request)

			err := authorizationHandler.Login(writer, request, responseFormat)
			if err != nil {
				log.Println("Something went wrong during authentication", err)
			} else {
				log.Println("Login successful")
			}
		} else {
			defaultHandler(writer, request)
		}
	}

	return f, nil
}

func getUploadHandlerFunc(config *Config, defaultHandler http.HandlerFunc) (http.HandlerFunc, error) {
	authorizationHandler, err := GetAuthorizationHandler(config)
	if err != nil {
		return nil, err
	}
	uploadHandler, err := GetUploadHandler(config)
	if err != nil {
		return nil, err
	}
	clientErrorHandler, err := GetClientErrorHandler()
	if err != nil {
		return nil, err
	}
	internalServerErrorHandler, err := GetInternalServerErrorHandler()
	if err != nil {
		return nil, err
	}

	f := func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("gfs-version", GFSVersion)
		if request.Method == "POST" {
			responseFormat := getResponseFormat(request)

			err := authorizationHandler.CheckAuthenticated(request)
			if err != nil {
				clientErrorHandler.Handle(writer, err, responseFormat, http.StatusUnauthorized)
				return
			}

			err = uploadHandler.Handle(writer, request)
			if err != nil {
				if err == ErrNoUploadingUp {
					clientErrorHandler.Handle(writer, err, responseFormat, http.StatusBadRequest)
					return
				}
				internalServerErrorHandler.Handle(writer, err, responseFormat)
			}

		} else {
			defaultHandler(writer, request)
		}
	}

	return f, nil
}
