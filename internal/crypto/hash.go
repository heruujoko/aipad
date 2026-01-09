package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// SimilarityThreshold defines the minimum similarity ratio to consider content as duplicate
const SimilarityThreshold = 0.80

// GenerateHash creates a SHA256 hash of the normalized content
func GenerateHash(content string) string {
	// Normalize: trim whitespace and lowercase
	normalized := strings.ToLower(strings.TrimSpace(content))
	hash := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(hash[:])
}

// Normalize prepares content for comparison by trimming and lowercasing
func Normalize(content string) string {
	return strings.ToLower(strings.TrimSpace(content))
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

// LevenshteinDistance calculates the edit distance between two strings
func LevenshteinDistance(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	// Create matrix
	matrix := make([][]int, len(a)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(b)+1)
	}

	// Initialize first column
	for i := 0; i <= len(a); i++ {
		matrix[i][0] = i
	}

	// Initialize first row
	for j := 0; j <= len(b); j++ {
		matrix[0][j] = j
	}

	// Fill in the rest of the matrix
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(a)][len(b)]
}

// SimilarityRatio calculates the similarity between two strings (0.0 to 1.0)
func SimilarityRatio(a, b string) float64 {
	a = Normalize(a)
	b = Normalize(b)

	if a == b {
		return 1.0
	}

	maxLen := max(len(a), len(b))
	if maxLen == 0 {
		return 1.0
	}

	distance := LevenshteinDistance(a, b)
	return 1.0 - float64(distance)/float64(maxLen)
}

// IsSimilar checks if the new content is similar to any existing content
func IsSimilar(newContent string, existingContents []string, threshold float64) (bool, string, float64) {
	for _, existing := range existingContents {
		ratio := SimilarityRatio(newContent, existing)
		if ratio >= threshold {
			return true, existing, ratio
		}
	}
	return false, "", 0
}
