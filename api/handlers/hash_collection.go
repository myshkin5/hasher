package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

type hashStore interface {
	AddPassword(password string) uint64
	GetHash(requestId uint64) (hash string, err error)
}

const HashCollectionPattern = "/hash"

func NewHashCollectionFunc(store hashStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("only POST method accepted"))
			return
		}

		r.ParseForm()

		passwords, ok := r.Form["password"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("password not set"))
			return
		}

		count := store.AddPassword(passwords[0])

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, strconv.FormatUint(count, 10))
	}
}
