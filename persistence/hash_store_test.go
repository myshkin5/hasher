package persistence_test

import (
	"testing"
	"time"

	"github.com/myshkin5/hasher/persistence"
)

func TestHashStore(t *testing.T) {
	t.Run("AddPassword", func(t *testing.T) {
		t.Run("increments request id", func(t *testing.T) {
			// ASSEMBLE
			store := persistence.NewHashStore(time.Millisecond, func(input string) string { return input })

			// ACT
			requestId1 := store.AddPassword("pass1")
			requestId2 := store.AddPassword("pass2")

			// ASSERT
			if requestId1 != 1 {
				t.Error("First request did not have request id of 1")
			}
			if requestId2 != 2 {
				t.Error("Second request did not have request id of 2")
			}
		})
	})

	t.Run("AddPassword/GetHash", func(t *testing.T) {
		t.Run("returns expected hash", func(t *testing.T) {
			// ASSEMBLE
			store := persistence.NewHashStore(time.Millisecond, func(input string) string { return "hash1" })

			// ACT
			requestId := store.AddPassword("pass1")
			time.Sleep(5 * time.Millisecond)
			hash, err := store.GetHash(requestId)

			// ASSERT
			if err != nil {
				t.Error("Returned unexpected error")
			}
			if hash != "hash1" {
				t.Error("Did not return expected hash")
			}
		})
	})
}
