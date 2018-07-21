package handlers_test

import "net/http/httptest"

func newRecorder() *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	w.Code = 0
	return w
}

type mockHashStore struct {
	addPasswordCalled    uint
	addPasswordPassword  string
	addPasswordRequestId uint64

	getHashCalled    uint
	getHashRequestId uint64
	getHashHash      string
	getHashErr       error
}

func (m *mockHashStore) AddPassword(password string) uint64 {
	m.addPasswordCalled++
	m.addPasswordPassword = password
	return m.addPasswordRequestId
}

func (m *mockHashStore) GetHash(requestId uint64) (hash string, err error) {
	m.getHashCalled++
	m.getHashRequestId = requestId
	return m.getHashHash, m.getHashErr
}
