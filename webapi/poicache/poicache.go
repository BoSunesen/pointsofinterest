package poicache

import (
	"encoding/json"
	"fmt"
	"github.com/BoSunesen/pointsofinterest/webapi/factories"
	"github.com/BoSunesen/pointsofinterest/webapi/logging"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

//TODO Less pointy?
type PoiCache struct {
	data          *[]PoiData //TODO Use MemCache to store data?
	dataMutex     *sync.RWMutex
	refreshMutex  *sync.Mutex
	expires       *time.Time //TODO Hard expire as well as soft
	logger        logging.Logger
	clientFactory factories.ClientFactory
	refresher     factories.BackgroundWorker
}

func NewPoiCache(logger logging.Logger, clientFactory factories.ClientFactory, workerFactory factories.WorkerFactory) *PoiCache {
	poiData := make([]PoiData, 0)
	expires := time.Now().Add(-42 * time.Hour)
	cache := PoiCache{data: &poiData, dataMutex: &sync.RWMutex{}, refreshMutex: &sync.Mutex{}, expires: &expires, logger: logger, clientFactory: clientFactory}

	cache.refresher = workerFactory.CreateBackgroundWorker("RefreshIfNeeded", func(delayedContext context.Context) {
		//TODO cache variable bound to refresher via closure, is that okay?
		cache.RefreshIfNeeded(delayedContext)
	})

	return &cache
}

func (cache *PoiCache) BackgroundRefresh(r *http.Request) error {
	return cache.refresher.DoWork(r)
}

func (cache *PoiCache) RefreshIfNeeded(ctx context.Context) {
	//TODO Load test
	if cache.expires.Before(time.Now()) {
		cache.logger.Infof(ctx, "Cache expired at %v waiting for refresh lock", cache.expires)

		cache.refreshMutex.Lock()
		defer cache.refreshMutex.Unlock()

		if cache.expires.Before(time.Now()) {
			cache.logger.Infof(ctx, "Refreshing cache")
			err := cache.refresh(ctx)
			if err != nil {
				cache.logger.Infof(ctx, "Error while refreshing cache: %v", err)
			} else {
				cache.logger.Infof(ctx, "Cache refreshed. Expires again at %v", cache.expires)
			}
		} else {
			cache.logger.Infof(ctx, "Got refresh lock, but cache is no longer expired. Expires again at %v", cache.expires)
		}
	}
}

func (cache *PoiCache) refresh(ctx context.Context) error {
	poiBytes, err := cache.getRemoteData(ctx, 10)
	if err != nil {
		return err
	}

	poiData := make([]PoiData, 0)
	err = json.Unmarshal(poiBytes, &poiData)
	if err != nil {
		return err
	}

	cache.logger.Infof(ctx, "Waiting for cache write lock")

	cache.dataMutex.Lock()
	defer cache.dataMutex.Unlock()

	cache.logger.Infof(ctx, "Got cache write lock")

	cache.data = &poiData

	//TODO How often to refresh
	expires := time.Now().Add(20 * time.Second)
	cache.expires = &expires

	return nil
}

func (cache *PoiCache) getRemoteData(ctx context.Context, retries int) ([]byte, error) {
	const remoteUrl string = "https://data.sfgov.org/resource/6a9r-agq8.json?status=APPROVED"
	request, err := http.NewRequest(http.MethodGet, remoteUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/json;charset=utf-8")
	if openDataAppToken := os.Getenv("OpenDataAppToken"); openDataAppToken != "" {
		cache.logger.Infof(ctx, "Setting X-App-Token header")
		request.Header.Set("X-App-Token", openDataAppToken)
	}

	client := cache.clientFactory.CreateClient(ctx)
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("Received status code %v while trying to refresh cache. Response body: %v", res.StatusCode, string(bytes))
	}

	if res.StatusCode == 202 {
		if retries <= 0 {
			return nil, fmt.Errorf("No more retries while trying to refresh cache. Response body: %v", string(bytes))
		}
		return cache.getRemoteData(ctx, retries-1)
	}

	return bytes, nil
}

func (cache *PoiCache) ReadData() *[]PoiData {
	cache.dataMutex.RLock()
	defer cache.dataMutex.RUnlock()

	poiData := make([]PoiData, len(*cache.data))
	copy(poiData, *cache.data)
	return &poiData
}
