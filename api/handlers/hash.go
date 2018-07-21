package handlers

import (
	"net/http"
	"strconv"

	"github.com/myshkin5/hasher/persistence"
)

const HashPattern = "/hash/"

func NewHashFunc(store hashStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("only GET method accepted"))
			return
		}

		requestId, err := strconv.ParseUint(r.URL.Path[len(HashPattern):], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad request id"))
			return
		}

		hash, err := store.GetHash(requestId)
		if err == persistence.ErrHashNotAvailable {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(hash))
	}
}
