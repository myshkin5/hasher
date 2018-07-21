package hasher

import (
	"crypto/sha512"
	"encoding/base64"
)

func SHA512(input string) string {
	data := sha512.Sum512([]byte(input))
	return base64.StdEncoding.EncodeToString(data[:])
}
