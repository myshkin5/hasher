package hasher_test

import (
	"testing"

	"github.com/myshkin5/hasher/hasher"
)

func TestSHA512(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// ASSEMBLE
		// ACT
		hash := hasher.SHA512("angryMonkey")

		// ASSERT
		if hash != "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==" {
			t.Error("hasher did not return expected hash")
		}
	})

	t.Run("emtpy string", func(t *testing.T) {
		// ASSEMBLE
		// ACT
		hash := hasher.SHA512("")

		// ASSERT
		if hash != "z4PhNX7vuL3xVChQ1m2AB9Yg5AULVxXcg/SpIdNs6c5H0NE8XYXysP+DGNKHfuwvY7kxvUdBeoGlODJ6+SfaPg==" {
			t.Error("hasher did not return expected hash")
		}
	})
}
