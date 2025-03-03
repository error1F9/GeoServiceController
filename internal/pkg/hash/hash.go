package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func Body(body []byte) string {
	hash := sha256.Sum256(body)
	return hex.EncodeToString(hash[:])
}
