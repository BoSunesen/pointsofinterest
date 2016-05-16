package factories

import (
	"golang.org/x/net/context"
	"net/http"
)

type ClientFactory interface {
	CreateClient(context.Context) *http.Client
}

type ContextFactory interface {
	CreateContext(r *http.Request) context.Context
}

type BackgroundWorker interface {
	DoWork(ctx context.Context) error
}

type WorkerFactory interface {
	CreateBackgroundWorker(string, func(context.Context)) BackgroundWorker
}
