package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/myshkin5/hasher/api/handlers"
)

func TestHashCollectionFunc(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// ARRANGE
		w := newRecorder()

		data := url.Values{}
		data.Set("password", "my-pass")
		r := httptest.NewRequest(http.MethodPost, "/hash", strings.NewReader(data.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		mockStore := mockHashStore{}
		mockStore.addPasswordRequestId = 42

		handlerFunc := handlers.NewHashCollectionFunc(&mockStore)

		// ACT
		handlerFunc.ServeHTTP(w, r)

		// ASSERT
		if mockStore.addPasswordPassword != "my-pass" {
			t.Error("handler did not add password")
		}
		if w.Code != http.StatusCreated {
			t.Error("handler did not return Created")
		}
		if w.Body.String() != "42" {
			t.Error("handler did not return request id body")
		}
	})

	t.Run("no password - bad request", func(t *testing.T) {
		// ARRANGE
		w := newRecorder()

		data := url.Values{}
		r := httptest.NewRequest(http.MethodPost, "/hash", strings.NewReader(data.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		mockStore := mockHashStore{}

		handlerFunc := handlers.NewHashCollectionFunc(&mockStore)

		// ACT
		handlerFunc.ServeHTTP(w, r)

		// ASSERT
		if mockStore.addPasswordCalled != 0 {
			t.Error("handler added empty password")
		}
		if w.Code != http.StatusBadRequest {
			t.Error("handler did not return Bad Request")
		}
	})

	t.Run("not POST - method not allowed", func(t *testing.T) {
		// ARRANGE
		w := newRecorder()
		r := httptest.NewRequest(http.MethodGet, "/hash", nil)

		mockStore := mockHashStore{}

		handlerFunc := handlers.NewHashCollectionFunc(&mockStore)

		// ACT
		handlerFunc.ServeHTTP(w, r)

		// ASSERT
		if mockStore.addPasswordCalled != 0 {
			t.Error("handler added empty password on GET")
		}
		if w.Code != http.StatusMethodNotAllowed {
			t.Error("handler did not return Method Not Allowed")
		}
	})
}
