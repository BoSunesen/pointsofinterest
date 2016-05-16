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
	DoWork(r *http.Request) error
}

type WorkerFactory interface {
	CreateBackgroundWorker(string, func(context.Context)) BackgroundWorker
}
