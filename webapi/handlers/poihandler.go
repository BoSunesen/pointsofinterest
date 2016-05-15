package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/BoSunesen/pointsofinterest/webapi/logging"
	"github.com/BoSunesen/pointsofinterest/webapi/poicache"
	"net/http"
)

type PoiHandler struct {
	cache  *poicache.PoiCache
	logger logging.Logger
}

func NewPoiHandler(logger logging.Logger) *PoiHandler {
	cache := poicache.NewPoiCache(logger)
	return &PoiHandler{cache, logger}
}

func (handler *PoiHandler) ServeHttp(w http.ResponseWriter, r *http.Request) error {
	go handler.cache.RefreshIfNeeded(r)

	poiData := handler.cache.ReadData()

	//TODO Only return POI's matching input (type, location, opening hours)
	jsonBytes, err := json.Marshal(poiData)
	if err != nil {
		return err
	}

	fmt.Fprint(w, string(jsonBytes))
	return nil
}
