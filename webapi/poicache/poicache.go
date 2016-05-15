package poicache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

//TODO Less pointy?
type PoiCache struct {
	data         *[]PoiData
	dataMutex    *sync.RWMutex
	refreshMutex *sync.Mutex
	expires      *time.Time
}

func NewPoiCache() *PoiCache {
	poiData := make([]PoiData, 0)
	expires := time.Now().Add(-42 * time.Hour)
	cache := PoiCache{data: &poiData, dataMutex: &sync.RWMutex{}, refreshMutex: &sync.Mutex{}, expires: &expires}
	cache.RefreshIfNeeded()
	return &cache
}

func (cache *PoiCache) RefreshIfNeeded() {
	//TODO Load test
	if cache.expires.Before(time.Now()) {
		log.Printf("Cache expired at %v waiting for refresh lock", cache.expires)

		cache.refreshMutex.Lock()
		defer cache.refreshMutex.Unlock()

		if cache.expires.Before(time.Now()) {
			log.Println("Refreshing cache")
			err := cache.refresh()
			if err != nil {
				log.Printf("Error while refreshing cache: %v", err)
			} else {
				log.Printf("Cache refreshed. Expires again at %v", cache.expires)
			}
		} else {
			log.Printf("Got refresh lock, but cache is no longer expired. Expires again at %v", cache.expires)
		}
	}
}

func (cache *PoiCache) refresh() error {
	poiBytes, err := getRemoteData(10)
	if err != nil {
		return err
	}

	poiData := make([]PoiData, 0)
	err = json.Unmarshal(poiBytes, &poiData)
	if err != nil {
		return err
	}

	log.Println("Waiting for cache write lock")

	cache.dataMutex.Lock()
	defer cache.dataMutex.Unlock()

	log.Println("Got cache write lock")

	cache.data = &poiData

	//TODO How often to refresh
	expires := time.Now().Add(1 * time.Minute)
	cache.expires = &expires

	return nil
}

func getRemoteData(retries int) ([]byte, error) {
	//TODO API key
	const remoteUrl string = "https://data.sfgov.org/resource/6a9r-agq8.json?status=APPROVED"
	request, err := http.NewRequest(http.MethodGet, remoteUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/json;charset=utf-8")
	if openDataAppToken := os.Getenv("OpenDataAppToken"); openDataAppToken != "" {
		request.Header.Set("X-App-Token", openDataAppToken)
	}

	client := http.Client{}
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
		return getRemoteData(retries - 1)
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
