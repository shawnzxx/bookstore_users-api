package crypto_utils

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMd5(input string) string {
	hash := md5.New()
	// Reset resets the Hash to its initial state.
	defer hash.Reset()
	hash.Write([]byte(input))
	// Sum appends the current hash to b and returns the resulting slice.
	// It does not change the underlying hash state.
	return hex.EncodeToString(hash.Sum(nil))
}
