package factories

import (
	"golang.org/x/net/context"
	"net/http"
)

type ClientFactory interface {
	// Create a client used to call the POI data provider.
	CreateClient(context.Context) *http.Client
}

type ContextFactory interface {
	// Create a context that is passed to the other factories
	// and to github.com/BoSunesen/pointsofinterest/webapi/logging.Logger
	//
	// This might be called more than once per request.
	CreateContext(r *http.Request) context.Context
}

type BackgroundWorker interface {
	// Perform work in the background, outside the scope of a user request. This should be a non-blocking call.
	DoWork(ctx context.Context) error
}

type WorkerFactory interface {
	// Create a BackgroundWorker that can execute the given function outside the scope of a user request.
	CreateBackgroundWorker(string, func(context.Context)) BackgroundWorker
}
