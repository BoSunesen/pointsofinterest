package initialization

import (
	"github.com/BoSunesen/pointsofinterest/webapi/factories"
	"github.com/BoSunesen/pointsofinterest/webapi/handlers"
	"github.com/BoSunesen/pointsofinterest/webapi/logging"
	"html"
	"net/http"
)

type WebApiInitializer struct {
	logger         logging.Logger
	contextFactory factories.ContextFactory
	clientFactory  factories.ClientFactory
	workerFactory  factories.WorkerFactory
}

func NewWebApiInitializer(
	logger logging.Logger,
	contextFactory factories.ContextFactory,
	clientFactory factories.ClientFactory,
	workerFactory factories.WorkerFactory,
) WebApiInitializer {
	return WebApiInitializer{
		logger,
		contextFactory,
		clientFactory,
		workerFactory,
	}
}

func (i WebApiInitializer) Initialize() {
	handle("/poi/", handlers.NewPoiHandler(i.logger, i.contextFactory, i.clientFactory, i.workerFactory), i.logger, i.contextFactory)

	const pingRoute string = "/ping/"
	handle(pingRoute, handlers.PingHandler{}, i.logger, i.contextFactory)

	//TODO Better routing
	//TODO favicon.ico
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := i.contextFactory.CreateContext(r)
		path := html.EscapeString(r.URL.Path)
		i.logger.Debugf(ctx, "Redirecting %v", path)
		defer i.logger.Debugf(ctx, "Redirected %v", path)
		http.Redirect(w, r, pingRoute, http.StatusFound)
	})
}

func handle(route string, handler handlers.HttpHandler, logger logging.Logger, contextFactory factories.ContextFactory) {
	http.Handle(route, &handlers.LoggingDecorator{handler, route, logger, contextFactory})
}
