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

	t.Run("GetHash", func(t *testing.T) {
		t.Run("unavailable hashes returns error", func(t *testing.T) {
			// ASSEMBLE
			store := persistence.NewHashStore(time.Millisecond, func(input string) string { return input })

			store.AddPassword("pass1")
			store.AddPassword("pass2")

			// ACT
			hash, err := store.GetHash(3)

			// ASSERT
			if err != persistence.ErrHashNotAvailable {
				t.Error("Error not returned for non-existent hash")
			}
			if hash != "" {
				t.Error("Non-empty hash returned for non-existent hash")
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

		t.Run("delays hash availability", func(t *testing.T) {
			// ASSEMBLE
			store := persistence.NewHashStore(500*time.Millisecond, func(input string) string { return "hash1" })

			// ACT
			// ASSERT
			requestId := store.AddPassword("pass1")
			start := time.Now()
			var hashAvailableTooSoon bool
			for {
				if time.Now().After(start.Add(450 * time.Millisecond)) {
					break
				}

				_, err := store.GetHash(requestId)
				if err != persistence.ErrHashNotAvailable {
					hashAvailableTooSoon = true
					break
				}

				time.Sleep(10 * time.Millisecond)
			}

			if hashAvailableTooSoon {
				t.Error("Hash available too soon")
			}

			time.Sleep(time.Now().Sub(start))

			hash, err := store.GetHash(requestId)
			if err != nil {
				t.Error("Returned unexpected error")
			}
			if hash != "hash1" {
				t.Error("Did not return expected hash")
			}
		})
	})
}
