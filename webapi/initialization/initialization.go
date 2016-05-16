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
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		ctx := i.contextFactory.CreateContext(r)
		path := html.EscapeString(r.URL.Path)
		i.logger.Debugf(ctx, "Serving favicon.ico, path: %v", path)
		defer i.logger.Debugf(ctx, "Served favicon.ico, path: %v", path)
		http.ServeFile(w, r, "favicon.ico")
	})

	poiHandler := handlers.NewPoiHandler(i.logger, i.contextFactory, i.clientFactory, i.workerFactory)
	handle("/poi", poiHandler, i.logger, i.contextFactory)
	handle("/poi/", poiHandler, i.logger, i.contextFactory)
	handle("/ping", handlers.PingHandler{}, i.logger, i.contextFactory)
}

func handle(route string, handler handlers.HttpHandler, logger logging.Logger, contextFactory factories.ContextFactory) {
	http.Handle(route, &handlers.LoggingDecorator{handler, route, logger, contextFactory})
}
