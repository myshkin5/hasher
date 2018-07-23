package integeration

import (
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/myshkin5/hasher/logs"
)

func StartRequesters(wg *sync.WaitGroup, requesterCount int, requestCount int) {
	for i := 0; i < requesterCount; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < requestCount; j++ {
				data := url.Values{}
				data.Set("password", "my-pass")
				response, err := http.Post(
					"http://localhost:8080/hash",
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
