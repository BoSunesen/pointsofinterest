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

type PoiCache struct {
	data          *[]PoiData
	dataMutex     *sync.RWMutex
	refreshMutex  *sync.Mutex
	staleAt       *time.Time
	refreshAt     *time.Time
	logger        logging.Logger
	clientFactory factories.ClientFactory
}

type PoiCacheRefresher struct {
	refreshWorker factories.BackgroundWorker
	cache         *PoiCache
}

func NewPoiCache(logger logging.Logger, clientFactory factories.ClientFactory, workerFactory factories.WorkerFactory) (*PoiCache, *PoiCacheRefresher) {
	poiData := make([]PoiData, 0)
	staleAt := time.Now().Add(-42 * time.Hour)
	refreshAt := time.Now().Add(-42 * time.Hour)
	cache := &PoiCache{
		&poiData,
		&sync.RWMutex{},
		&sync.Mutex{},
		&staleAt,
		&refreshAt,
		logger,
		clientFactory,
	}

	refresher := workerFactory.CreateBackgroundWorker("backgroundRefresh", func(delayedContext context.Context) {
		cache.backgroundRefresh(delayedContext)
	})

	return cache, &PoiCacheRefresher{refresher, cache}
}

func (refresher *PoiCacheRefresher) Refresh(ctx context.Context) error {
	if refresher.cache.IsStale() {
		return refresher.cache.immediateRefresh(ctx)
	}
	if refresher.cache.NeedsRefresh() {
		return refresher.refreshWorker.DoWork(ctx)
	}
	return nil
}

func (cache *PoiCache) IsStale() bool {
	return cache.staleAt.Before(time.Now())
}

func (cache *PoiCache) NeedsRefresh() bool {
	return cache.refreshAt.Before(time.Now())
}

func (cache *PoiCache) immediateRefresh(ctx context.Context) error {
	cache.logger.Warningf(ctx, "Cache went stale at %v, waiting for refresh lock", cache.staleAt)

	cache.refreshMutex.Lock()
	defer cache.refreshMutex.Unlock()

	if cache.IsStale() {
		cache.logger.Warningf(ctx, "Refreshing stale cache")
		err := cache.refresh(ctx)
		if err != nil {
			return err
		} else {
			cache.logger.Warningf(ctx, "Stale cache refreshed. Refresh again at %v, stale at %v", cache.refreshAt, cache.staleAt)
		}
	} else {
		cache.logger.Warningf(ctx, "Got refresh lock, but cache no longer stale. Refresh again at %v, stale at %v", cache.refreshAt, cache.staleAt)
	}
	return nil
}

func (cache *PoiCache) backgroundRefresh(ctx context.Context) {
	cache.logger.Infof(ctx, "Cache needed refresh at %v, waiting for refresh lock", cache.refreshAt)

	cache.refreshMutex.Lock()
	defer cache.refreshMutex.Unlock()

	if cache.NeedsRefresh() {
		cache.logger.Infof(ctx, "Refreshing cache")
		err := cache.refresh(ctx)
		if err != nil {
			cache.logger.Infof(ctx, "Error while refreshing cache: %v", err)
		} else {
			cache.logger.Infof(ctx, "Cache refreshed. Refresh again at %v", cache.refreshAt)
		}
	} else {
		cache.logger.Infof(ctx, "Got refresh lock, but cache no longer needs refresh. Refresh again at %v", cache.refreshAt)
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

	staleAt := time.Now().Add(4 * time.Hour)
	cache.staleAt = &staleAt

	refreshAt := time.Now().Add(30 * time.Second)
	cache.refreshAt = &refreshAt

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
