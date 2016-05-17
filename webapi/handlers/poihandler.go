package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/BoSunesen/pointsofinterest/webapi/factories"
	"github.com/BoSunesen/pointsofinterest/webapi/logging"
	"github.com/BoSunesen/pointsofinterest/webapi/poicache"
	"github.com/BoSunesen/pointsofinterest/webapi/poidata"
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

	//TODO Cache parsed data
	parser := poidata.PoiParser{handler.logger}
	parsedData := parser.ParsePoiData(ctx, *poiData)

	//TODO Only return POI's matching input (type, location, opening hours)
	jsonBytes, err := json.Marshal(parsedData)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	fmt.Fprint(w, string(jsonBytes))
	return nil
}
