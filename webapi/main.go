package main

import (
	"github.com/BoSunesen/pointsofinterest/webapi/initialization"
	"log"
	"net/http"
)

func main() {
	initialization.Initialize()
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
