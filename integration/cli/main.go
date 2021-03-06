package main

import (
	"os"
	"strconv"
	"sync"

	"github.com/myshkin5/hasher/integration"
	"github.com/myshkin5/hasher/logs"
)

func main() {
	err := logs.Init(getEnvWithDefault("LOG_LEVEL", "info"))
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		logs.Logger.Panic("First CLI option must be defined")
	}

	wg := &sync.WaitGroup{}

	serverURL := getEnvWithDefault("SERVER_URL", "http://localhost:8080")

	switch os.Args[1] {
	case "requesters":
		startRequesters(wg, serverURL)
	default:
		logs.Logger.Panic("First CLI option must be one of: requesters")
	}

	wg.Wait()
}

func startRequesters(wg *sync.WaitGroup, serverURL string) {
	if len(os.Args) != 4 {
		logs.Logger.Panic("CLI: <executable> requesters <requesterCount> <requestCount>")
	}

	requesterCount, err := strconv.Atoi(os.Args[2])
	if err != nil {
		logs.Logger.Panicf("Could not get requester count: %s", err.Error())
	}
	requestCount, err := strconv.Atoi(os.Args[3])
	if err != nil {
		logs.Logger.Panicf("Could not get request count: %s", err.Error())
	}

	wg.Add(requesterCount)
	integration.StartRequesters(wg, requesterCount, requestCount, serverURL)
}

func getEnvWithDefault(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}

	return defaultValue
}
