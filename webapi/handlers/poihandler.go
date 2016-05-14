package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type PoiHandler struct {
	cache *poiCache
}

func NewPoiHandler() *PoiHandler {
	poiData := make([]poiData, 0)
	expires := time.Now().Add(-42 * time.Hour)
	cache := &poiCache{data: &poiData, dataMutex: &sync.RWMutex{}, refreshMutex: &sync.Mutex{}, expires: &expires}
	cache.refreshIfNeeded()
	return &PoiHandler{cache}
}

func (handler *PoiHandler) ServeHttp(w http.ResponseWriter, r *http.Request) error {
	go handler.cache.refreshIfNeeded()

	jsonBytes, err := handler.readCache()
	if err != nil {
		return err
	}

	fmt.Fprint(w, string(jsonBytes))
	return nil
}

func (handler *PoiHandler) readCache() ([]byte, error) {
	handler.cache.dataMutex.RLock()
	defer handler.cache.dataMutex.RUnlock()

	//TODO Only return POI's matching input (type, location, opening hours)

	return json.Marshal(handler.cache.data)
}
