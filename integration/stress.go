package integration

import (
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/myshkin5/hasher/logs"
)

func StartRequesters(wg *sync.WaitGroup, requesterCount int, requestCount int, serverURL string) {
	for i := 0; i < requesterCount; i++ {
		go func() {
			defer wg.Done()
			client := &http.Client{}
			for j := 0; j < requestCount; j++ {
				data := url.Values{}
				data.Set("password", "my-pass")
				response, err := client.Post(
					serverURL+"/hash",
					"application/x-www-form-urlencoded",
					strings.NewReader(data.Encode()))
				if err != nil {
					logs.Logger.Errorf("Could not request hash: %s", err.Error())
				}
				if response != nil && response.StatusCode != http.StatusCreated {
					logs.Logger.Errorf("Did not receive Created on requesting hash: %d", response.StatusCode)
				}
				if response != nil && response.Body != nil {
					response.Body.Close()
				}
			}
		}()
	}
}
