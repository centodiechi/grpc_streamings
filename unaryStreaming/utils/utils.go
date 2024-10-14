package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(key string) string {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
