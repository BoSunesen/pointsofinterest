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

type GoWorkerFactory struct {
	ContextFactory factories.ContextFactory
}

func (factory GoWorkerFactory) CreateBackgroundWorker(key string, workFunction func(context.Context)) factories.BackgroundWorker {
	return GoWorker{workFunction, factory.ContextFactory}
}

type GoWorker struct {
	workFunction   func(context.Context)
	contextFactory factories.ContextFactory
}

func (w GoWorker) DoWork(r *http.Request) error {
	ctx := w.contextFactory.CreateContext(r)
	go w.workFunction(ctx)
	return nil
}
