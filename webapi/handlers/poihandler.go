package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/BoSunesen/pointsofinterest/webapi/poicache"
	"net/http"
)

type PoiHandler struct {
	cache *poicache.PoiCache
}

func NewPoiHandler() *PoiHandler {
	cache := poicache.NewPoiCache()
	return &PoiHandler{cache}
}

func (handler *PoiHandler) ServeHttp(w http.ResponseWriter, r *http.Request) error {
	go handler.cache.RefreshIfNeeded()

	poiData := handler.cache.ReadData()

	//TODO Only return POI's matching input (type, location, opening hours)
	jsonBytes, err := json.Marshal(poiData)
	if err != nil {
		return err
	}

	fmt.Fprint(w, string(jsonBytes))
	return nil
}
