package persistence

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/myshkin5/hasher/logs"
	"github.com/myshkin5/hasher/metrics"
)

type HashFunc func(string) string

var ErrHashNotAvailable = errors.New("store: the requested hash is not available")

type HashStore struct {
	delay         time.Duration
	hashFunc      HashFunc
	storeCount    uint
	hashStopwatch *metrics.Stopwatch

	count uint64

	hashes []atomic.Value
}

type hash struct {
	requestId uint64
	available time.Time
	hash      string
}

func NewHashStore(delay time.Duration, hashFunc HashFunc, storeCount uint, hashStopwatch *metrics.Stopwatch) *HashStore {
	return &HashStore{
		delay:         delay,
		hashFunc:      hashFunc,
		storeCount:    storeCount,
		hashStopwatch: hashStopwatch,

		hashes: make([]atomic.Value, storeCount),
	}
}

func (s *HashStore) AddPassword(password string) uint64 {
	requestId := atomic.AddUint64(&s.count, 1)
	i := s.ringIndex(requestId)
	logs.Logger.Infof("Adding password, request %d/index %d...", requestId, i)

	start := time.Now()
	h := s.hashFunc(password)
	s.hashStopwatch.AddRun(start, time.Now())

	s.hashes[i].Store(hash{
		requestId: requestId,
		available: time.Now().Add(s.delay),
		hash:      h})
	return requestId
}

func (s *HashStore) GetHash(requestId uint64) (string, error) {
	i := s.ringIndex(requestId)

	hash, ok := s.hashes[i].Load().(hash)
	if !ok {
		logs.Logger.Infof("Hash not available for request %d/index %d, nil value", requestId, i)
		return "", ErrHashNotAvailable
	}

	if hash.requestId != requestId {
		logs.Logger.Infof("Hash not available for request %d/index %d, hash overwritten?", requestId, i)
		return "", ErrHashNotAvailable
	}

	if time.Now().Before(hash.available) {
		logs.Logger.Infof("Hash not available for request %d/index %d, insufficient delay", requestId, i)
		return "", ErrHashNotAvailable
	}

	logs.Logger.Infof("Returning hash for request %d/index %d", requestId, i)

	return hash.hash, nil
}

func (s *HashStore) ringIndex(requestId uint64) int {
	return int((requestId - 1) % uint64(s.storeCount))
}
