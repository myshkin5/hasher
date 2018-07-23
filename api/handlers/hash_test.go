package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/myshkin5/hasher/api/handlers"
	"github.com/myshkin5/hasher/persistence"
)

func TestHashFunc(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// ARRANGE
		w := newRecorder()
		r := httptest.NewRequest(http.MethodGet, "/hash/42", nil)

		mockStore := mockHashStore{}
		mockStore.getHashHash = "hash-123"

		handlerFunc := handlers.NewHashFunc(&mockStore)

		// ACT
		handlerFunc.ServeHTTP(w, r)

		// ASSERT
		if mockStore.getHashRequestId != 42 {
			t.Error("handler did not pass right request id")
		}
		if w.Code != http.StatusOK {
			t.Error("handler did not return OK")
		}
		if w.Body.String() != "hash-123" {
			t.Error("handler did not return hash body")
		}
	})

	t.Run("not GET - method not allowed", func(t *testing.T) {
		// ARRANGE
		w := newRecorder()
		r := httptest.NewRequest(http.MethodPost, "/hash/42", nil)

		mockStore := mockHashStore{}

		handlerFunc := handlers.NewHashFunc(&mockStore)

		// ACT
		handlerFunc.ServeHTTP(w, r)

		// ASSERT
		if mockStore.getHashCalled != 0 {
			t.Error("handler got hash on POST")
		}
		if w.Code != http.StatusMethodNotAllowed {
			t.Error("handler did not return Method Not Allowed")
		}
	})

	t.Run("empty request id - bad request", func(t *testing.T) {
		// ARRANGE
		w := newRecorder()
		r := httptest.NewRequest(http.MethodGet, "/hash/", nil)

		mockStore := mockHashStore{}
		mockStore.getHashHash = "hash-123"

		handlerFunc := handlers.NewHashFunc(&mockStore)

		// ACT
		handlerFunc.ServeHTTP(w, r)

		// ASSERT
		if mockStore.getHashCalled != 0 {
			t.Error("handler got hash with no request id")
		}
		if w.Code != http.StatusBadRequest {
			t.Error("handler did not return Bad Request")
		}
	})

	t.Run("non-numeric request id - bad request", func(t *testing.T) {
		// ARRANGE
		w := newRecorder()
		r := httptest.NewRequest(http.MethodGet, "/hash/whatcha", nil)

		mockStore := mockHashStore{}
		mockStore.getHashHash = "hash-123"

		handlerFunc := handlers.NewHashFunc(&mockStore)

		// ACT
		handlerFunc.ServeHTTP(w, r)

		// ASSERT
		if mockStore.getHashCalled != 0 {
			t.Error("handler got hash with no request id")
		}
		if w.Code != http.StatusBadRequest {
			t.Error("handler did not return Bad Request")
		}
	})

	t.Run("hash not available - not found", func(t *testing.T) {
		// ARRANGE
		w := newRecorder()
		r := httptest.NewRequest(http.MethodGet, "/hash/42", nil)

		mockStore := mockHashStore{}
		mockStore.getHashErr = persistence.ErrHashNotAvailable

		handlerFunc := handlers.NewHashFunc(&mockStore)

		// ACT
		handlerFunc.ServeHTTP(w, r)

		// ASSERT
		if w.Code != http.StatusNotFound {
			t.Error("handler did not return Not Found")
		}
	})

	t.Run("other store error - internal server error", func(t *testing.T) {
		// ARRANGE
		w := newRecorder()
		r := httptest.NewRequest(http.MethodGet, "/hash/42", nil)

		mockStore := mockHashStore{}
		mockStore.getHashErr = errors.New("off the rails")

		handlerFunc := handlers.NewHashFunc(&mockStore)

		// ACT
		handlerFunc.ServeHTTP(w, r)

		// ASSERT
		if w.Code != http.StatusInternalServerError {
			t.Error("handler did not return internal server error")
		}
	})
}
