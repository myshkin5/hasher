package persistence_test

import (
	"sync"
	"testing"
	"time"

	"github.com/myshkin5/hasher/persistence"
)

func TestHashStore(t *testing.T) {
	t.Run("AddPassword", func(t *testing.T) {
		t.Run("increments request id", func(t *testing.T) {
			// ASSEMBLE
			store := persistence.NewHashStore(time.Millisecond, func(input string) string { return input }, 1)

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

		t.Run("avoids race conditions", func(t *testing.T) {
			// ASSEMBLE
			store := persistence.NewHashStore(time.Millisecond, func(input string) string { return "hash1" }, 1)

			start := time.Now()
			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					if time.Now().After(start.Add(450 * time.Millisecond)) {
						break
					}

					store.AddPassword("pass")

					time.Sleep(10 * time.Millisecond)
				}
			}()

			// ACT
			// ASSERT (expect race detector will assert any issues)
			store.AddPassword("pass")

			wg.Wait()
		})
	})

	t.Run("GetHash", func(t *testing.T) {
		t.Run("unavailable hashes returns error", func(t *testing.T) {
			// ASSEMBLE
			store := persistence.NewHashStore(time.Millisecond, func(input string) string { return input }, 3)

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
			store := persistence.NewHashStore(time.Millisecond, func(input string) string { return "hash1" }, 1)

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
			store := persistence.NewHashStore(500*time.Millisecond, func(input string) string { return "hash1" }, 1)

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

		t.Run("avoids race conditions", func(t *testing.T) {
			// ASSEMBLE
			store := persistence.NewHashStore(500*time.Millisecond, func(input string) string { return "hash1" }, 1)

			start := time.Now()
			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					if time.Now().After(start.Add(500 * time.Millisecond)) {
						break
					}

					store.GetHash(1)

					time.Sleep(10 * time.Millisecond)
				}
			}()

			// ACT
			// ASSERT (expect race detector will assert any issues)
			store.AddPassword("pass")

			wg.Wait()
		})

		t.Run("returns right value in ring buffer", func(t *testing.T) {
			// ASSEMBLE
			store := persistence.NewHashStore(time.Millisecond, func(input string) string { return input }, 2)

			store.AddPassword("hash1")
			store.AddPassword("hash2")
			store.AddPassword("hash3")

			time.Sleep(5 * time.Millisecond)

			// ACT
			_, err1 := store.GetHash(1)
			hash2, err2 := store.GetHash(2)
			hash3, err3 := store.GetHash(3)

			// ASSERT
			if err1 != persistence.ErrHashNotAvailable {
				t.Error("Returned unexpected value")
			}
			if err2 != nil {
				t.Error("Returned unexpected error")
			}
			if hash2 != "hash2" {
				t.Error("Returned unexpected hash")
			}
			if err3 != nil {
				t.Error("Returned unexpected error")
			}
			if hash3 != "hash3" {
				t.Error("Returned unexpected hash")
			}
		})
	})
}
