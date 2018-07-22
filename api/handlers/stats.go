package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/myshkin5/hasher/logs"
	"github.com/myshkin5/hasher/metrics"
)

const StatsPattern = "/stats"

func NewStatsFunc(stopwatch *metrics.Stopwatch) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("only GET method accepted"))
			return
		}

		stats := stopwatch.Statistics()
		logs.Logger.Infof("Got statistics: %#v", stats)

		s := struct {
			Total   uint64 `json:"total"`
			Average int64  `json:"average"`
		}{Total: stats.Total, Average: stats.Average.Nanoseconds() / 1000}
		out, _ := json.Marshal(s)

		w.WriteHeader(http.StatusOK)
		w.Write(out)
	}
}
