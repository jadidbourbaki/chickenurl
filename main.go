package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
)

var urlToShortMap = make(map[string]string)
var shortToUrlMap = make(map[string]string)
var urlMutex = sync.Mutex{}

func init() {
	// Load production quality logger
	logger := zap.Must(zap.NewProduction())
	if os.Getenv("APP_ENV") == "development" {
		logger = zap.Must(zap.NewDevelopment())
	}

	zap.ReplaceGlobals(logger)
}

func main() {
	http.HandleFunc("/", mainHandler)
	go saveMapsToFile()

	http.ListenAndServe(":8080", nil)
}

func writeMapToGob(theMap map[string]string, filename string) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	urlMutex.Lock()
	enc.Encode(theMap)
	urlMutex.Unlock()

	os.WriteFile(filename, buf.Bytes(), 0644)
	zap.L().Info("written map to file: " + filename)
}

func saveMapsToFile() {
	for {
		writeMapToGob(urlToShortMap, "urlToShortMap.gob")
		writeMapToGob(shortToUrlMap, "shortToUrlMap.gob")

		time.Sleep(1 * time.Second)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postHandler(w, r)
		return
	}

	if r.Method == http.MethodDelete {
		deleteHandler(w, r)
		return
	}

	if r.Method == http.MethodGet {
		getHandler(w, r)
		return
	}

	zap.L().Error("invalid http method: " + r.Method)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	shortKey, err := io.ReadAll(r.Body)
	if err != nil {
		zap.L().Error(err.Error())
		return
	}

	stringifiedShortKey := string(shortKey)

	urlMutex.Lock()
	gotUrl, ok := shortToUrlMap[stringifiedShortKey]
	urlMutex.Unlock()

	if !ok {
		fmt.Fprintf(w, "short key not fund %v", stringifiedShortKey)
		return
	}

	http.Redirect(w, r, gotUrl, http.StatusMovedPermanently)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	url, err := io.ReadAll(r.Body)
	if err != nil {
		zap.L().Error(err.Error())
		return
	}

	stringifiedUrl := string(url)

	urlMutex.Lock()
	gotShortKey, ok := urlToShortMap[stringifiedUrl]
	urlMutex.Unlock()

	if ok {
		fmt.Fprintf(w, "url %v shortened to %v", stringifiedUrl, gotShortKey)
		return
	}

	shortKey := urlToShort()

	urlMutex.Lock()
	shortToUrlMap[shortKey] = stringifiedUrl
	urlToShortMap[stringifiedUrl] = shortKey
	urlMutex.Unlock()

	fmt.Fprintf(w, "url %v shortened to %v", stringifiedUrl, shortKey)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	shortKey, err := io.ReadAll(r.Body)
	if err != nil {
		zap.L().Error(err.Error())
		return
	}

	stringifiedShortKey := string(shortKey)

	urlMutex.Lock()
	gotUrl, ok := shortToUrlMap[stringifiedShortKey]
	urlMutex.Unlock()

	if !ok {
		fmt.Fprintf(w, "short url %v not found", stringifiedShortKey)
		return
	}

	urlMutex.Lock()
	delete(shortToUrlMap, stringifiedShortKey)
	delete(urlToShortMap, gotUrl)
	urlMutex.Unlock()

	fmt.Fprintf(w, "short url %v deleted", stringifiedShortKey)
}
