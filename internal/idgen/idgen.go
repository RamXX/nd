package idgen

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand/v2"
	"time"
)

// GenerateID creates a collision-resistant ID in the form PREFIX-HASH.
// The hash is 3 hex chars derived from SHA-256 of title + timestamp + nonce.
// existsFn is called to check for collisions; it retries with a new nonce if needed.
func GenerateID(prefix, title string, existsFn func(string) bool) (string, error) {
	const maxRetries = 100
	for i := range maxRetries {
		nonce := fmt.Sprintf("%d-%d", time.Now().UnixNano(), i)
		raw := title + nonce
		hash := sha256.Sum256([]byte(raw))
		short := hex.EncodeToString(hash[:])[:3]
		id := fmt.Sprintf("%s-%s", prefix, short)
		if existsFn == nil || !existsFn(id) {
			return id, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique ID after %d attempts", maxRetries)
}

// GenerateChildID creates a child ID in the form PARENT.N where N is sequential.
// nextN should be the next available child number (1-based).
func GenerateChildID(parentID string, nextN int) string {
	return fmt.Sprintf("%s.%d", parentID, nextN)
}

// RandomNonce returns a short random hex string for extra entropy.
func RandomNonce() string {
	b := make([]byte, 4)
	for i := range b {
		b[i] = byte(rand.IntN(256))
	}
	return hex.EncodeToString(b)
}
