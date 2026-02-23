package enforce

import (
	"crypto/sha256"
	"fmt"
)

// ComputeContentHash returns the SHA-256 hash of the body content.
// The hash covers only the markdown body (below frontmatter), not frontmatter fields.
func ComputeContentHash(body string) string {
	h := sha256.Sum256([]byte(body))
	return fmt.Sprintf("sha256:%x", h)
}
