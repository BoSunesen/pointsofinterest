package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BoSunesen/pointsofinterest/webapi/factories"
	"github.com/BoSunesen/pointsofinterest/webapi/logging"
	"github.com/BoSunesen/pointsofinterest/webapi/poicache"
	"github.com/BoSunesen/pointsofinterest/webapi/poidata"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
	"strconv"
)

type PoiHandler struct {
	cache          *poicache.PoiCache
	refresher      *poicache.PoiCacheRefresher
	logger         logging.Logger
	contextFactory factories.ContextFactory
}

func NewPoiHandler(logger logging.Logger, contextFactory factories.ContextFactory, clientFactory factories.ClientFactory, workerFactory factories.WorkerFactory) *PoiHandler {
	cache, refresher := poicache.InitCache(logger, clientFactory, workerFactory)
	return &PoiHandler{cache, refresher, logger, contextFactory}
}

func (handler *PoiHandler) ServeHttp(w http.ResponseWriter, r *http.Request) error {
	ctx := handler.contextFactory.CreateContext(r)
	err := handler.refresher.Refresh(ctx)
	if err != nil {
		errMessage := fmt.Sprintf("Error while refreshing cache: %v", err)
		if handler.cache.IsStale() {
			return errors.New(errMessage)
		}
	}

	queryValues := r.URL.Query()

	latitude, latitudeIsParsed, err := handler.parseFloatInput(ctx, queryValues, "latitude")
	if err != nil {
		return err
	}
	longitude, longitudeIsParsed, err := handler.parseFloatInput(ctx, queryValues, "longitude")
	if err != nil {
		return err
	}
	distance, distanceIsParsed, err := handler.parseIntInput(ctx, queryValues, "distance")
	if err != nil {
		return err
	}
	weekday, weekdayIsParsed, err := handler.parseIntInput(ctx, queryValues, "weekday")
	if err != nil {
		return err
	}
	hour, hourIsParsed, err := handler.parseIntInput(ctx, queryValues, "hour")
	if err != nil {
		return err
	}

	applyGeographicFilter := latitudeIsParsed && longitudeIsParsed && distanceIsParsed
	applyOpeningHoursFilter := weekdayIsParsed || hourIsParsed

	rawData := handler.cache.ReadData()

	parser := poidata.PoiParser{handler.logger}
	output := parser.ParsePoiData(ctx, rawData)

	if applyGeographicFilter {
		output = poidata.FilterByLocation(output, latitude, longitude, distance)
	}

	if applyOpeningHoursFilter {
		output = poidata.FilterByOpeningHours(output, weekday, hour)
	}

	jsonBytes, err := json.MarshalIndent(*output, "", " ")
	if err != nil {
		return err
	}

	handler.logger.Debugf(ctx, "Returning %v elements", len(*output))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, string(jsonBytes))
	return nil
}

func (h PoiHandler) parseFloatInput(ctx context.Context, queryValues url.Values, key string) (float64, bool, error) {
	valueString := queryValues.Get(key)
	if len(valueString) > 0 {
		value, err := strconv.ParseFloat(valueString, 64)
		if err != nil {
			return -1, false, fmt.Errorf("Error while parsing %v: %v", key, err)
		}
		return value, true, nil
	}
	return -1, false, nil
}

func (h PoiHandler) parseIntInput(ctx context.Context, queryValues url.Values, key string) (int, bool, error) {
	valueString := queryValues.Get(key)
	if len(valueString) > 0 {
		value, err := strconv.Atoi(valueString)
		if err != nil {
			return -1, false, fmt.Errorf("Error while parsing %v: %v", key, err)
		}
		return value, true, nil
	}
	return -1, false, nil
}
