package handlers

import (
	"fmt"
	"github.com/BoSunesen/pointsofinterest/webapi/factories"
	"github.com/BoSunesen/pointsofinterest/webapi/logging"
	"html"
	"net/http"
	"runtime/debug"
)

type HttpHandler interface {
	ServeHttp(http.ResponseWriter, *http.Request) error
}

type LoggingDecorator struct {
	Handler        HttpHandler
	Route          string
	Logger         logging.Logger
	ContextFactory factories.ContextFactory
}

func (decorator *LoggingDecorator) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	ctx := decorator.ContextFactory.CreateContext(request)
	path := html.EscapeString(request.URL.Path)
	decorator.Logger.Debugf(ctx, "Serving %v", path)
	defer decorator.Logger.Debugf(ctx, "Served %v", path)

	defer func() {
		if r := recover(); r != nil {
			errorString := fmt.Sprintln(r, string(debug.Stack()))
			decorator.Logger.Criticalf(ctx, "Panic while handling route %v (path: %v):\n%v", decorator.Route, path, errorString)
			http.Error(w, errorString, http.StatusInternalServerError)
		}
	}()

	err := decorator.Handler.ServeHttp(w, request)

	if err != nil {
		errorString := err.Error()
		decorator.Logger.Errorf(ctx, "Error while handling route %v (path: %v): %v", decorator.Route, path, errorString)
		http.Error(w, errorString, http.StatusInternalServerError)
	}
}
