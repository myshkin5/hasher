package metrics_test

import (
	"testing"
	"time"

	"github.com/myshkin5/hasher/metrics"
)

func TestStopwatch(t *testing.T) {
	t.Run("AddRun/Statistics", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			// ASSEMBLE
			now := time.Now()
			stopwatch := metrics.Stopwatch{}
			stopwatch.AddRun(now, now.Add(time.Second))
			stopwatch.AddRun(now.Add(3*time.Second), now.Add(5*time.Second))

			// ACT
			statistics := stopwatch.Statistics()

			// ASSERT
			if statistics.Average != 1500*time.Millisecond {
				t.Error("Unexpected average")
			}
			if statistics.Total != 2 {
				t.Error("Unexpected total")
			}
		})

		t.Run("doesn't divide by zero when no runs are added", func(t *testing.T) {
			// ASSEMBLE
			stopwatch := metrics.Stopwatch{}

			// ACT
			statistics := stopwatch.Statistics()

			// ASSERT
			if statistics.Average != 0 {
				t.Error("Unexpected average")
			}
			if statistics.Total != 0 {
				t.Error("Unexpected total")
			}
		})
	})
}
