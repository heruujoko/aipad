package crypto

import (
	"testing"
)

func TestGenerateHash(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "simple text",
			input: "hello world",
			// SHA256 of "hello world" (already lowercase)
			expected: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		},
		{
			name:     "text with whitespace",
			input:    "  hello world  ",
			expected: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		},
		{
			name:     "text with uppercase",
			input:    "HELLO WORLD",
			expected: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateHash(tt.input)
			if result != tt.expected {
				t.Errorf("GenerateHash(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsDuplicate(t *testing.T) {
	existingHashes := []string{"hash1", "hash2", "hash3"}

	tests := []struct {
		name     string
		hash     string
		expected bool
	}{
		{"exists", "hash2", true},
		{"not exists", "hash4", false},
		{"empty hash", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDuplicate(tt.hash, existingHashes)
			if result != tt.expected {
				t.Errorf("IsDuplicate(%q) = %v, want %v", tt.hash, result, tt.expected)
			}
		})
	}
}

func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		a, b     string
		expected int
	}{
		{"", "", 0},
		{"abc", "", 3},
		{"", "abc", 3},
		{"abc", "abc", 0},
		{"abc", "abd", 1},
		{"kitten", "sitting", 3},
		{"saturday", "sunday", 3},
	}

	for _, tt := range tests {
		t.Run(tt.a+"_"+tt.b, func(t *testing.T) {
			result := LevenshteinDistance(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("LevenshteinDistance(%q, %q) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestSimilarityRatio(t *testing.T) {
	tests := []struct {
		a, b        string
		minExpected float64
		maxExpected float64
	}{
		{"hello", "hello", 1.0, 1.0},
		{"hello", "HELLO", 1.0, 1.0}, // case insensitive
		{"hello", "hallo", 0.7, 0.9},
		{"hello world", "hello world!", 0.9, 1.0},
		{"abc", "xyz", 0.0, 0.1},
	}

	for _, tt := range tests {
		t.Run(tt.a+"_"+tt.b, func(t *testing.T) {
			result := SimilarityRatio(tt.a, tt.b)
			if result < tt.minExpected || result > tt.maxExpected {
				t.Errorf("SimilarityRatio(%q, %q) = %f, want between %f and %f", tt.a, tt.b, result, tt.minExpected, tt.maxExpected)
			}
		})
	}
}

func TestIsSimilar(t *testing.T) {
	existingContents := []string{
		"Discussed the new API design with focus on REST principles",
		"Implemented user authentication with JWT tokens",
	}

	tests := []struct {
		name          string
		newContent    string
		threshold     float64
		expectSimilar bool
	}{
		{
			name:          "exact match",
			newContent:    "Discussed the new API design with focus on REST principles",
			threshold:     0.8,
			expectSimilar: true,
		},
		{
			name:          "slight variation",
			newContent:    "Discussed the new API design with focus on REST principle",
			threshold:     0.8,
			expectSimilar: true,
		},
		{
			name:          "completely different",
			newContent:    "Refactored database schema for better performance",
			threshold:     0.8,
			expectSimilar: false,
		},
		{
			name:          "case variation",
			newContent:    "DISCUSSED THE NEW API DESIGN WITH FOCUS ON REST PRINCIPLES",
			threshold:     0.8,
			expectSimilar: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			similar, _, _ := IsSimilar(tt.newContent, existingContents, tt.threshold)
			if similar != tt.expectSimilar {
				t.Errorf("IsSimilar(%q) = %v, want %v", tt.newContent, similar, tt.expectSimilar)
			}
		})
	}
}
