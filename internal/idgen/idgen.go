package idgen

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"math/rand/v2"
	"strings"
	"time"
)

// base36Alphabet is the character set for base36 encoding (0-9, a-z).
const base36Alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"

// EncodeBase36 converts a byte slice to a base36 string of specified length.
func EncodeBase36(data []byte, length int) string {
	num := new(big.Int).SetBytes(data)

	var result strings.Builder
	base := big.NewInt(36)
	zero := big.NewInt(0)
	mod := new(big.Int)

	chars := make([]byte, 0, length)
	for num.Cmp(zero) > 0 {
		num.DivMod(num, base, mod)
		chars = append(chars, base36Alphabet[mod.Int64()])
	}

	// Reverse
	for i := len(chars) - 1; i >= 0; i-- {
		result.WriteByte(chars[i])
	}

	str := result.String()
	if len(str) < length {
		str = strings.Repeat("0", length-len(str)) + str
	}
	if len(str) > length {
		str = str[len(str)-length:]
	}

	return str
}

// GenerateID creates a collision-resistant ID in the form PREFIX-HASH.
// The hash is 4 base36 chars derived from SHA-256 of title + timestamp + nonce.
// existsFn is called to check for collisions; it retries with a new nonce if needed.
func GenerateID(prefix, title string, existsFn func(string) bool) (string, error) {
	const maxRetries = 100
	for i := range maxRetries {
		nonce := fmt.Sprintf("%d-%d", time.Now().UnixNano(), i)
		raw := title + nonce
		hash := sha256.Sum256([]byte(raw))
		short := EncodeBase36(hash[:3], 4) // 3 bytes -> 4 base36 chars
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
	var sb strings.Builder
	for _, c := range b {
		fmt.Fprintf(&sb, "%02x", c)
	}
	return sb.String()
}
