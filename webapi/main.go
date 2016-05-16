package main

import (
	"github.com/BoSunesen/pointsofinterest/webapi/initialization"
	"log"
	"net/http"
)

func main() {
	webApiInitializer := initialization.NewWebApiInitializer(GoLog{}, GoContextFactory{}, GoClientFactory{}, GoWorkerFactory{})
	webApiInitializer.Initialize()
	err := http.ListenAndServe(":8080", nil)
	log.Fatalln(err)
}
