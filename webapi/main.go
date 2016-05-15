package main

import (
	"github.com/BoSunesen/pointsofinterest/webapi/initialization"
	"github.com/BoSunesen/pointsofinterest/webapi/logging"
	"log"
	"net/http"
)

func main() {
	initialization.Initialize(logging.GoLog{})
	err := http.ListenAndServe(":8080", nil)
	log.Fatalln(err)
}
