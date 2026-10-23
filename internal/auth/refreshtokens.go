package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	key := make([]byte, 128)
	rand.Read(key)
	keyStr := hex.EncodeToString(key)
	return keyStr
}
