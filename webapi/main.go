package main

import (
	"github.com/BoSunesen/pointsofinterest/webapi/initialization"
	"log"
	"net/http"
)

func main() {
	cf := GoContextFactory{}
	webApiInitializer := initialization.NewWebApiInitializer(GoLog{}, cf, GoClientFactory{}, GoWorkerFactory{cf})
	webApiInitializer.Initialize()
	err := http.ListenAndServe(":8080", nil)
	log.Fatalln(err)
}
