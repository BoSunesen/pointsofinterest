package main

import (
	"github.com/BoSunesen/pointsofinterest/webapi/initialization"
	"log"
	"net/http"
)

func main() {
	cf := GoContextFactory{}
	initialization.Initialize(GoLog{}, cf, GoClientFactory{}, GoWorkerFactory{cf})
	err := http.ListenAndServe(":8080", nil)
	log.Fatalln(err)
}
