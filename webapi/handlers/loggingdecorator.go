package handlers

import (
	"fmt"
	"github.com/BoSunesen/pointsofinterest/webapi/logging"
	"html"
	"net/http"
)

type HttpHandler interface {
	ServeHttp(http.ResponseWriter, *http.Request) error
}

type LoggingDecorator struct {
	Handler HttpHandler
	Route   string
	//TODO LoggerFactory?
	Logger  logging.Logger
}

func (decorator *LoggingDecorator) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	path := html.EscapeString(request.URL.Path)
	decorator.Logger.Debugf(request, "Serving %v", path)
	defer decorator.Logger.Debugf(request, "Served %v", path)

	defer func() {
		if r := recover(); r != nil {
			errorString := fmt.Sprint(r)
			decorator.Logger.Criticalf(request, "Panic while handling route %v (path: %v): %v", decorator.Route, path, errorString)
			http.Error(w, errorString, http.StatusInternalServerError)
		}
	}()

	err := decorator.Handler.ServeHttp(w, request)

	if err != nil {
		errorString := err.Error()
		decorator.Logger.Errorf(request, "Error while handling route %v (path: %v): %v", decorator.Route, path, errorString)
		http.Error(w, errorString, http.StatusInternalServerError)
	}
}
