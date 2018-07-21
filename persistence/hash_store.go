package persistence

import (
	"errors"
	"time"

	"github.com/myshkin5/hasher/logs"
)

var ErrHashNotAvailable = errors.New("store: the requested hash is not available")

type HashStore struct {
}

func NewHashStore(delay time.Duration) *HashStore {
	return &HashStore{}
}

func (s *HashStore) AddPassword(password string) uint64 {
	logs.Logger.Info("Added password ...")
	return 0
}

func (s *HashStore) GetHash(requestId uint64) (string, error) {
	logs.Logger.Infof("Returning hash for request %d", requestId)
	return "hash", nil
}

//func (h *Hasher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	c := atomic.AddUint64(&h.count, 1)
//	fmt.Fprint(w, strconv.FormatUint(c, 10))
//}
