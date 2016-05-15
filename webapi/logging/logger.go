package logging

import (
	"log"
	"net/http"
)

type Logger interface {
	Debugf(*http.Request, string, ...interface{})
	Infof(*http.Request, string, ...interface{})
	Warningf(*http.Request, string, ...interface{})
	Errorf(*http.Request, string, ...interface{})
	Criticalf(*http.Request, string, ...interface{})
}

type GoLog struct{}

func (l GoLog) Debugf(r *http.Request, format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l GoLog) Infof(r *http.Request, format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l GoLog) Warningf(r *http.Request, format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l GoLog) Errorf(r *http.Request, format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l GoLog) Criticalf(r *http.Request, format string, v ...interface{}) {
	log.Printf(format, v...)
}
