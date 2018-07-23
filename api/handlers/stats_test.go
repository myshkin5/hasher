package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/myshkin5/hasher/api/handlers"
	"github.com/myshkin5/hasher/metrics"
)

func TestStatsFunc(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// ARRANGE
		w := newRecorder()
		r := httptest.NewRequest(http.MethodGet, "/stats", nil)

		now := time.Now()
		stopwatch := metrics.Stopwatch{}
		stopwatch.AddRun(now, now.Add(time.Second))

		handlerFunc := handlers.NewStatsFunc(&stopwatch)

		// ACT
		handlerFunc.ServeHTTP(w, r)

		// ASSERT
		if w.Code != http.StatusOK {
			t.Error("Handler did not return OK")
		}
		if w.Body.String() != `{"total":1,"average":1000000}` {
			t.Error("Unexpected stats body")
		}
	})

	t.Run("not GET - method not allowed", func(t *testing.T) {
		// ARRANGE
		w := newRecorder()
		r := httptest.NewRequest(http.MethodPost, "/stats", nil)

		stopwatch := metrics.Stopwatch{}
		handlerFunc := handlers.NewStatsFunc(&stopwatch)

		// ACT
		handlerFunc.ServeHTTP(w, r)

		// ASSERT
		if w.Code != http.StatusMethodNotAllowed {
			t.Error("handler did not return Method Not Allowed")
		}
	})
}
