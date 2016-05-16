package main

import (
	"golang.org/x/net/context"
	"log"
)

type GoLog struct{}

func (l GoLog) Debugf(ctx context.Context, format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l GoLog) Infof(ctx context.Context, format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l GoLog) Warningf(ctx context.Context, format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l GoLog) Errorf(ctx context.Context, format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l GoLog) Criticalf(ctx context.Context, format string, v ...interface{}) {
	log.Printf(format, v...)
}
