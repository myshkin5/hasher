package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/myshkin5/hasher/api/handlers"
)

func TestShutdownFunc(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// ARRANGE
		w := newRecorder()
		r := httptest.NewRequest(http.MethodGet, "/shutdown", nil)

		shuttingDown := make(chan struct{}, 0)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-shuttingDown
		}()

		handlerFunc := handlers.NewShutdownFunc(shuttingDown)

		// ACT
		handlerFunc.ServeHTTP(w, r)

		// ASSERT (will hang forever if shuttingDown channel isn't closed)
		wg.Wait()
	})
}
