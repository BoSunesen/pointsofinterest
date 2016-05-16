package logging

import (
	//TODO depend on golang.org/x/net/context?
	"golang.org/x/net/context"
)

type Logger interface {
	Debugf(context.Context, string, ...interface{})
	Infof(context.Context, string, ...interface{})
	Warningf(context.Context, string, ...interface{})
	Errorf(context.Context, string, ...interface{})
	Criticalf(context.Context, string, ...interface{})
}