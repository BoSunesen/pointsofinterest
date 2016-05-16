package main

import (
	"github.com/BoSunesen/pointsofinterest/webapi/factories"
	"golang.org/x/net/context"
	"net/http"
)

type GoClientFactory struct{}

func (f GoClientFactory) CreateClient(ctx context.Context) *http.Client {
	return http.DefaultClient
}

type GoContextFactory struct{}

func (f GoContextFactory) CreateContext(r *http.Request) context.Context {
	return context.TODO()
}

type GoWorkerFactory struct{}

func (factory GoWorkerFactory) CreateBackgroundWorker(key string, workFunction func(context.Context)) factories.BackgroundWorker {
	return GoWorker{workFunction}
}

type GoWorker struct {
	workFunction func(context.Context)
}

func (w GoWorker) DoWork(ctx context.Context) error {
	go w.workFunction(ctx)
	return nil
}
