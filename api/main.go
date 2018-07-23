package main

import (
	"net/http"
	"os"
	"time"

	"github.com/myshkin5/hasher/api/handlers"
	"github.com/myshkin5/hasher/hash"
	"github.com/myshkin5/hasher/logs"
	"github.com/myshkin5/hasher/metrics"
	"github.com/myshkin5/hasher/persistence"
)

func main() {
	initLogging()

	hashStopwatch := metrics.Stopwatch{}

	store := persistence.NewHashStore(5*time.Second, hash.SHA512, 10000, &hashStopwatch)

	mux := initRoutes(store, &hashStopwatch)

	serverAddr := getEnvWithDefault("SERVER_ADDR", "localhost")
	port := getEnvWithDefault("PORT", "8080")

	listenAndServe(serverAddr, port, mux)
}

func initLogging() {
	err := logs.Init(getEnvWithDefault("LOG_LEVEL", "info"))
	if err != nil {
		panic(err)
	}
}

func initRoutes(store *persistence.HashStore, hashStopwatch *metrics.Stopwatch) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc(handlers.HashCollectionPattern, handlers.NewHashCollectionFunc(store))
	mux.HandleFunc(handlers.HashPattern, handlers.NewHashFunc(store))
	mux.HandleFunc(handlers.StatsPattern, handlers.NewStatsFunc(hashStopwatch))

	return mux
}

func listenAndServe(serverAddr, port string, mux *http.ServeMux) {
	logs.Logger.Infof("Listening on %s:%s...", serverAddr, port)
	err := http.ListenAndServe(serverAddr+":"+port, mux)
	if err != nil {
		logs.Logger.Panic("ListenAndServe: ", err)
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}

	return defaultValue
}
