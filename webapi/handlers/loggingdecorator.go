package handlers

import (
	"log"
	"net/http"
)

type AppHandler interface {
	ServeHttpInner(http.ResponseWriter, *http.Request) error
}

type LoggingDecorator struct {
	Handler AppHandler
	AppName string
}

func (decorator LoggingDecorator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := decorator.Handler.ServeHttpInner(w, r)
	if err != nil {
		var errorString = err.Error()
		log.Printf("Error from %s: %q", decorator.AppName, errorString)
		http.Error(w, errorString, http.StatusInternalServerError)
	}
}
