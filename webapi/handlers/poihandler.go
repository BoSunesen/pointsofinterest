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
	refresher      *poicache.PoiCacheRefresher
	logger         logging.Logger
	contextFactory factories.ContextFactory
}

func NewPoiHandler(logger logging.Logger, contextFactory factories.ContextFactory, clientFactory factories.ClientFactory, workerFactory factories.WorkerFactory) *PoiHandler {
	cache, refresher := poicache.NewPoiCache(logger, clientFactory, workerFactory)
	return &PoiHandler{cache, refresher, logger, contextFactory}
}

func (handler *PoiHandler) ServeHttp(w http.ResponseWriter, r *http.Request) error {
	ctx := handler.contextFactory.CreateContext(r)
	err := handler.refresher.Refresh(ctx)
	if err != nil {
		handler.logger.Errorf(ctx, "Error while refreshing cache: %v", err)
		if handler.cache.IsStale() {
			return err
		}
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
