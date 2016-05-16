package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/BoSunesen/pointsofinterest/webapi/factories"
	"github.com/BoSunesen/pointsofinterest/webapi/logging"
	"github.com/BoSunesen/pointsofinterest/webapi/poicache"
	"net/http"
)

type PoiHandler struct {
	cache          *poicache.PoiCache
	logger         logging.Logger
	contextFactory factories.ContextFactory
}

func NewPoiHandler(logger logging.Logger, contextFactory factories.ContextFactory, clientFactory factories.ClientFactory, workerFactory factories.WorkerFactory) *PoiHandler {
	cache := poicache.NewPoiCache(logger, clientFactory, workerFactory)
	return &PoiHandler{cache, logger, contextFactory}
}

func (handler *PoiHandler) ServeHttp(w http.ResponseWriter, r *http.Request) error {
	err := handler.cache.BackgroundRefresh(r)
	if err != nil {
		handler.logger.Errorf(handler.contextFactory.CreateContext(r), "Could not start background refresh of cache: %v", err)
	}

	poiData := handler.cache.ReadData()

	//TODO Only return POI's matching input (type, location, opening hours)
	jsonBytes, err := json.Marshal(poiData)
	if err != nil {
		return err
	}

	fmt.Fprint(w, string(jsonBytes))
	return nil
}
