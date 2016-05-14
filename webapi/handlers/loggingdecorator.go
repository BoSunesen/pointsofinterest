package handlers

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

type HttpHandler interface {
	ServeHttp(http.ResponseWriter, *http.Request) error
}

type LoggingDecorator struct {
	Handler HttpHandler
	Route   string
}

func (decorator *LoggingDecorator) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	path := html.EscapeString(request.URL.Path)
	log.Printf("Serving %q", path)
	defer log.Printf("Served %q", path)

	defer func() {
		if r := recover(); r != nil {
			errorString := fmt.Sprint(r)
			log.Printf("Panic while handling route %q (path: %s): %q", decorator.Route, path, errorString)
			http.Error(w, errorString, http.StatusInternalServerError)
		}
	}()

	err := decorator.Handler.ServeHttp(w, request)

	if err != nil {
		errorString := err.Error()
		log.Printf("Error while handling route %q (path: %s): %q", decorator.Route, path, errorString)
		http.Error(w, errorString, http.StatusInternalServerError)
	}
}
