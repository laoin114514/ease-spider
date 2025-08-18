package component

import (
	"crypto/sha512"
	"encoding/hex"
)

func Hash(data string) string {
	byteData := []byte(data)
	hash := sha512.Sum512(byteData)
	hashHex := hex.EncodeToString(hash[:])
	return hashHex
}
