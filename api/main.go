package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/myshkin5/hasher/api/handlers"
	"github.com/myshkin5/hasher/hash"
	"github.com/myshkin5/hasher/logs"
	"github.com/myshkin5/hasher/metrics"
	"github.com/myshkin5/hasher/persistence"
)

func main() {
	initLogging()

	shuttingDown := make(chan struct{}, 0)
	hashStopwatch := &metrics.Stopwatch{}

	hashStoreCount, err := strconv.Atoi(getEnvWithDefault("HASH_STORE_COUNT", "10000"))
	if err != nil {
		logs.Logger.Panic("Bad value set for HASH_STORE_COUNT environment variable: ", err)
	}

	go trackHashRate(hashStopwatch, shuttingDown)

	store := persistence.NewHashStore(5*time.Second, hash.SHA512, uint(hashStoreCount), hashStopwatch)

	mux := initRoutes(store, hashStopwatch, shuttingDown)

	serverAddr := getEnvWithDefault("SERVER_ADDR", "localhost")
	port := getEnvWithDefault("PORT", "8080")

	server := listenAndServe(serverAddr, port, mux)

	<-shuttingDown

	server.Shutdown(context.Background())

	logs.Logger.Info("Shutdown complete.")
}

func initLogging() {
	err := logs.Init(getEnvWithDefault("LOG_LEVEL", "info"))
	if err != nil {
		panic(err)
	}
}

func trackHashRate(hashStopwatch *metrics.Stopwatch, shuttingDown chan struct{}) {
	ticker := time.NewTicker(time.Minute)
	var previousTotal uint64
	for {
		select {
		case <-ticker.C:
			newTotal := hashStopwatch.Statistics().Total
			logs.Logger.Warnf("Processed %d hashes in the last minute", newTotal-previousTotal)
			previousTotal = newTotal
		case <-shuttingDown:
			break
		}
	}
}

func initRoutes(store *persistence.HashStore, hashStopwatch *metrics.Stopwatch, shuttingDown chan struct{}) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc(handlers.HashCollectionPattern, handlers.NewHashCollectionFunc(store))
	mux.HandleFunc(handlers.HashPattern, handlers.NewHashFunc(store))
	mux.HandleFunc(handlers.StatsPattern, handlers.NewStatsFunc(hashStopwatch))
	mux.HandleFunc(handlers.ShutdownPattern, handlers.NewShutdownFunc(shuttingDown))

	return mux
}

func listenAndServe(serverAddr, port string, mux *http.ServeMux) *http.Server {
	logs.Logger.Infof("Listening on %s:%s...", serverAddr, port)
	server := http.Server{Addr: serverAddr + ":" + port, Handler: mux}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logs.Logger.Panic("Error returned by ListenAndServe: ", err)
		}
	}()

	return &server
}

func getEnvWithDefault(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}

	return defaultValue
}
