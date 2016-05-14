package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

//TODO Move poiCache to own module
type poiCache struct {
	data         *[]poiData
	dataMutex    *sync.RWMutex
	refreshMutex *sync.Mutex
	expires      *time.Time
}

func (cache *poiCache) refreshIfNeeded() {
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
				log.Println("Cache refreshed")
			}
		} else {
			log.Printf("Got refresh lock, but cache is no longer expired. Expires again at %q", cache.expires)
		}
	}
}

func (cache *poiCache) refresh() error {
	poiBytes, err := getRemoteData()

	poiData := make([]poiData, 0)
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
	if res.Body != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	//TODO Check http status before reading body
	return ioutil.ReadAll(res.Body)
}
