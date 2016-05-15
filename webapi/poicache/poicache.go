package poicache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

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
		log.Printf("Cache expired at %q waiting for refresh lock", cache.expires)

		cache.refreshMutex.Lock()
		defer cache.refreshMutex.Unlock()

		if cache.expires.Before(time.Now()) {
			log.Println("Refreshing cache")
			err := cache.refresh()
			if err != nil {
				log.Printf("Error while refreshing cache: %q", err)
			} else {
				log.Printf("Cache refreshed. Expires again at %q", cache.expires)
			}
		} else {
			log.Printf("Got refresh lock, but cache is no longer expired. Expires again at %q", cache.expires)
		}
	}
}

func (cache *PoiCache) refresh() error {
	poiBytes, err := getRemoteData()

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

func getRemoteData() ([]byte, error) {
	//TODO API key
	const remoteUrl string = "https://data.sfgov.org/resource/6a9r-agq8.json?status=APPROVED"
	request, err := http.NewRequest(http.MethodGet, remoteUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/json;charset=utf-8")
	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 100 || res.StatusCode >= 300 {
		errorString := fmt.Sprintf("Received status code %s while trying to refresh cache", res.StatusCode)
		log.Println(errorString)
		return nil, errors.New(errorString)
	}

	return ioutil.ReadAll(res.Body)
}

func (cache *PoiCache) ReadData() *[]PoiData {
	cache.dataMutex.RLock()
	defer cache.dataMutex.RUnlock()

	poiData := make([]PoiData, len(*cache.data))
	copy(poiData, *cache.data)
	return &poiData
}
