package metrics

import (
	"sync/atomic"
	"time"
)

type Stopwatch struct {
	total    uint64
	duration uint64
}

type Statistics struct {
	Average time.Duration
	Total   uint64
}

func (s *Stopwatch) AddRun(start time.Time, end time.Time) {
	atomic.AddUint64(&s.total, 1)
	atomic.AddUint64(&s.duration, uint64(end.Sub(start)))
}

func (s *Stopwatch) Statistics() Statistics {
	total := atomic.LoadUint64(&s.total)
	var average uint64
	if total > 0 {
		average = atomic.LoadUint64(&s.duration) / total
	}
	return Statistics{Average: time.Duration(average), Total: total}
}
