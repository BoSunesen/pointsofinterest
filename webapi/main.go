package main

import (
	"github.com/BoSunesen/pointsofinterest/webapi/handlers"
	"log"
	"net/http"
)

func main() {
	http.Handle("/ping/", handlers.LoggingDecorator{handlers.PingHandler{}, "Ping"})
	http.Handle("/poi/", handlers.LoggingDecorator{handlers.PoiHandler{}, "POI"})
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
