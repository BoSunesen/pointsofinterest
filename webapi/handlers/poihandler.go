package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type PoiHandler struct {
}

type PoiData struct {
	Latitude  string
	Longitude string
	Title     string
}

func (handler PoiHandler) ServeHttpInner(w http.ResponseWriter, r *http.Request) error {
	poiBytes, err := json.Marshal(PoiData{"123", "123", "Dummy POI"})
	if err != nil {
		return err
	}
	poiString := string(poiBytes)
	log.Println(poiString)
	fmt.Fprint(w, poiString)
	return nil
}
