package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// GenerateHash creates a SHA256 hash of the normalized content
func GenerateHash(content string) string {
	// Normalize: trim whitespace and lowercase
	normalized := strings.ToLower(strings.TrimSpace(content))
	hash := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(hash[:])
}

// IsDuplicate checks if the hash exists in the list of existing hashes
func IsDuplicate(hash string, existingHashes []string) bool {
	for _, h := range existingHashes {
		if h == hash {
			return true
		}
	}
	return false
}
