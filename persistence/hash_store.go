package persistence

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/myshkin5/hasher/logs"
)

type HashFunc func(string) string

var ErrHashNotAvailable = errors.New("store: the requested hash is not available")

type HashStore struct {
	delay    time.Duration
	hashFunc HashFunc

	count uint64

	hashes []hash
}

type hash struct {
	available time.Time
	hash      string
}

func NewHashStore(delay time.Duration, hashFunc HashFunc) *HashStore {
	return &HashStore{
		delay:    delay,
		hashFunc: hashFunc,

		hashes: make([]hash, 0),
	}
}

func (s *HashStore) AddPassword(password string) uint64 {
	c := atomic.AddUint64(&s.count, 1)
	logs.Logger.Infof("Adding password, request %d...", c)
	s.hashes = append(s.hashes, hash{available: time.Now().Add(s.delay), hash: s.hashFunc(password)})
	return c
}

func (s *HashStore) GetHash(requestId uint64) (string, error) {
	logs.Logger.Infof("Returning hash for request %d", requestId)

	if len(s.hashes) < int(requestId) {
		return "", ErrHashNotAvailable
	}

	hash := s.hashes[requestId-1]

	if time.Now().Before(hash.available) {
		return "", ErrHashNotAvailable
	}

	return hash.hash, nil
}
