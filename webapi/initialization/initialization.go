package initialization

import (
	"github.com/BoSunesen/pointsofinterest/webapi/handlers"
	"github.com/BoSunesen/pointsofinterest/webapi/logging"
	"html"
	"net/http"
)

func Initialize(logger logging.Logger) {
	handle("/poi/", handlers.NewPoiHandler(logger), logger)

	const pingRoute string = "/ping/"
	handle(pingRoute, handlers.PingHandler{}, logger)

	//TODO Better routing
	//TODO favicon.ico
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := html.EscapeString(r.URL.Path)
		logger.Debugf(r, "Redirecting %v", path)
		defer logger.Debugf(r, "Redirected %v", path)
		http.Redirect(w, r, pingRoute, http.StatusFound)
	})
}

func handle(route string, handler handlers.HttpHandler, logger logging.Logger) {
	http.Handle(route, &handlers.LoggingDecorator{Handler: handler, Route: route, Logger: logger})
}
