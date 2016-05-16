package initialization

import (
	"github.com/BoSunesen/pointsofinterest/webapi/factories"
	"github.com/BoSunesen/pointsofinterest/webapi/handlers"
	"github.com/BoSunesen/pointsofinterest/webapi/logging"
	"html"
	"net/http"
)

//TODO Move parameters to struct
func Initialize(logger logging.Logger, contextFactory factories.ContextFactory, clientFactory factories.ClientFactory, workerFactory factories.WorkerFactory) {
	handle("/poi/", handlers.NewPoiHandler(logger, contextFactory, clientFactory, workerFactory), logger, contextFactory)

	const pingRoute string = "/ping/"
	handle(pingRoute, handlers.PingHandler{}, logger, contextFactory)

	//TODO Better routing
	//TODO favicon.ico
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := contextFactory.CreateContext(r)
		path := html.EscapeString(r.URL.Path)
		logger.Debugf(ctx, "Redirecting %v", path)
		defer logger.Debugf(ctx, "Redirected %v", path)
		http.Redirect(w, r, pingRoute, http.StatusFound)
	})
}

func handle(route string, handler handlers.HttpHandler, logger logging.Logger, contextFactory factories.ContextFactory) {
	http.Handle(route, &handlers.LoggingDecorator{handler, route, logger, contextFactory})
}
