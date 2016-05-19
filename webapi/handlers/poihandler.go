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
	"time"
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
	applyOpeningHoursFilter := weekdayIsParsed && hourIsParsed

	poiData := handler.cache.ReadData()

	//TODO Cache parsed data?
	parser := poidata.PoiParser{handler.logger}
	output := parser.ParsePoiData(ctx, poiData)

	if applyGeographicFilter {
		output = filterOutputGeographically(output, latitude, longitude, distance)
	}

	if applyOpeningHoursFilter {
		output, err = filterOutputByOpeningHours(output, weekday, hour)
		if err != nil {
			return err
		}
	}

	//TODO Only return POI's matching input (type, status)
	jsonBytes, err := json.Marshal(*output)
	if err != nil {
		return err
	}

	handler.logger.Debugf(ctx, "Returning %v elements", len(*output))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	fmt.Fprint(w, string(jsonBytes))
	return nil
}

func (h PoiHandler) parseFloatInput(ctx context.Context, queryValues url.Values, key string) (float64, bool, error) {
	valueString := queryValues.Get(key)
	if len(valueString) > 0 {
		value, err := strconv.ParseFloat(valueString, 64)
		if err != nil {
			return 0, false, fmt.Errorf("Error while parsing %v: %v", key, err)
		}
		return value, true, nil
	}
	return 0, false, nil
}

func (h PoiHandler) parseIntInput(ctx context.Context, queryValues url.Values, key string) (int, bool, error) {
	valueString := queryValues.Get(key)
	if len(valueString) > 0 {
		value, err := strconv.Atoi(valueString)
		if err != nil {
			return 0, false, fmt.Errorf("Error while parsing %v: %v", key, err)
		}
		return value, true, nil
	}
	return 0, false, nil
}

func filterOutputGeographically(parsedData *[]poidata.ParsedPoiData, latitude, longitude float64, distance int) *[]poidata.ParsedPoiData {
	geoFilteredData := make([]poidata.ParsedPoiData, 0, len(*parsedData))
	boundingBox := poidata.CalculateBoundingBox(latitude, longitude, distance)
	for _, v := range *parsedData {
		if v.Latitude <= boundingBox.MaxLatitude &&
			v.Latitude >= boundingBox.MinLatitude &&
			v.Longitude <= boundingBox.MaxLongitude &&
			v.Longitude >= boundingBox.MinLongitude {
			geoFilteredData = append(geoFilteredData, v)
		}
	}
	return &geoFilteredData
}

func filterOutputByOpeningHours(parsedData *[]poidata.ParsedPoiData, weekday, hour int) (*[]poidata.ParsedPoiData, error) {
	openingsFilteredData := make([]poidata.ParsedPoiData, 0, len(*parsedData))
	var openingHoursFunc func(*poidata.OpeningHours) []poidata.TimeInterval
	switch time.Weekday(weekday) {
	case time.Sunday:
		openingHoursFunc = func(openingHours *poidata.OpeningHours) []poidata.TimeInterval { return openingHours.Sunday }
	case time.Monday:
		openingHoursFunc = func(openingHours *poidata.OpeningHours) []poidata.TimeInterval { return openingHours.Monday }
	case time.Tuesday:
		openingHoursFunc = func(openingHours *poidata.OpeningHours) []poidata.TimeInterval { return openingHours.Tuesday }
	case time.Wednesday:
		openingHoursFunc = func(openingHours *poidata.OpeningHours) []poidata.TimeInterval { return openingHours.Wednesday }
	case time.Thursday:
		openingHoursFunc = func(openingHours *poidata.OpeningHours) []poidata.TimeInterval { return openingHours.Thursday }
	case time.Friday:
		openingHoursFunc = func(openingHours *poidata.OpeningHours) []poidata.TimeInterval { return openingHours.Friday }
	case time.Saturday:
		openingHoursFunc = func(openingHours *poidata.OpeningHours) []poidata.TimeInterval { return openingHours.Saturday }
	default:
		return nil, fmt.Errorf("Unknown weekday: %v", weekday)
	}
	for _, element := range *parsedData {
		includeElement := false
		if element.OpeningHours != nil {
			for _, interval := range openingHoursFunc(element.OpeningHours) {
				if interval.OpenFrom <= hour && interval.OpenTo > hour {
					includeElement = true
					break
				}
			}
		}
		if includeElement {
			openingsFilteredData = append(openingsFilteredData, element)
		}
	}
	return &openingsFilteredData, nil
}
