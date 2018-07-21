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
	hashFunc HashFunc

	count uint64

	hashes []string
}

func NewHashStore(delay time.Duration, hashFunc HashFunc) *HashStore {
	return &HashStore{
		hashFunc: hashFunc,
		hashes:   make([]string, 0),
	}
}

func (s *HashStore) AddPassword(password string) uint64 {
	c := atomic.AddUint64(&s.count, 1)
	logs.Logger.Infof("Adding password, request %d...", c)
	s.hashes = append(s.hashes, s.hashFunc(password))
	return c
}

func (s *HashStore) GetHash(requestId uint64) (string, error) {
	logs.Logger.Infof("Returning hash for request %d", requestId)
	return s.hashes[requestId-1], nil
}
