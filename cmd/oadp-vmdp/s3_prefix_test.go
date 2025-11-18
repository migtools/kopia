package main

import (
	"strings"
	"testing"
)

func TestNormalizeOADPPrefix(t *testing.T) {
	tests := []struct {
		name        string
		userPrefix  string
		expected    string
		expectError bool
	}{
		{
			name:        "prefix with trailing slash",
			userPrefix:  "abc/",
			expected:    "oadp-vmdp/abc/",
			expectError: false,
		},
		{
			name:        "prefix without trailing slash",
			userPrefix:  "abc",
			expected:    "oadp-vmdp/abc",
			expectError: false,
		},
		{
			name:        "prefix with leading and trailing slash",
			userPrefix:  "/abc/bce/",
			expected:    "oadp-vmdp/abc/bce/",
			expectError: false,
		},
		{
			name:        "prefix with leading slash no trailing slash",
			userPrefix:  "/abc/bce",
			expected:    "oadp-vmdp/abc/bce",
			expectError: false,
		},
		{
			name:        "just slash",
			userPrefix:  "/",
			expected:    "oadp-vmdp/",
			expectError: false,
		},
		{
			name:        "empty string",
			userPrefix:  "",
			expected:    "oadp-vmdp/",
			expectError: false,
		},
		{
			name:        "multiple leading slashes",
			userPrefix:  "///abc/",
			expected:    "oadp-vmdp/abc/",
			expectError: false,
		},
		// Error cases: user provides oadp-vmdp in prefix
		{
			name:        "user provides exact oadp-vmdp prefix",
			userPrefix:  "oadp-vmdp/",
			expected:    "",
			expectError: true,
		},
		{
			name:        "user provides oadp-vmdp with leading slash",
			userPrefix:  "/oadp-vmdp/",
			expected:    "",
			expectError: true,
		},
		{
			name:        "user provides oadp-vmdp in middle",
			userPrefix:  "some/oadp-vmdp/path",
			expected:    "",
			expectError: true,
		},
		{
			name:        "user provides uppercase OADP-VMDP",
			userPrefix:  "OADP-VMDP/",
			expected:    "",
			expectError: true,
		},
		{
			name:        "user provides mixed case OaDp-VmDp",
			userPrefix:  "OaDp-VmDp/backup",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NormalizeOADPPrefix(tt.userPrefix)
			if tt.expectError {
				if err == nil {
					t.Errorf("NormalizeOADPPrefix(%q) should return an error", tt.userPrefix)
				} else if !containsString(err.Error(), "must not contain 'oadp-vmdp'") {
					t.Errorf("Expected error to contain 'must not contain 'oadp-vmdp'', got: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("NormalizeOADPPrefix(%q) should not return an error: %v", tt.userPrefix, err)
				}
				if result != tt.expected {
					t.Errorf("NormalizeOADPPrefix(%q) = %q, want %q", tt.userPrefix, result, tt.expected)
				}
			}
		})
	}
}

func TestDenormalizeOADPPrefix(t *testing.T) {
	tests := []struct {
		name             string
		normalizedPrefix string
		expected         string
	}{
		{
			name:             "prefix with trailing slash",
			normalizedPrefix: "oadp-vmdp/abc/",
			expected:         "abc/",
		},
		{
			name:             "prefix without trailing slash",
			normalizedPrefix: "oadp-vmdp/abc",
			expected:         "abc",
		},
		{
			name:             "just oadp prefix",
			normalizedPrefix: "oadp-vmdp/",
			expected:         "",
		},
		{
			name:             "no oadp prefix",
			normalizedPrefix: "something-else/",
			expected:         "something-else/",
		},
		{
			name:             "empty string",
			normalizedPrefix: "",
			expected:         "",
		},
		{
			name:             "nested paths",
			normalizedPrefix: "oadp-vmdp/path/to/backup/",
			expected:         "path/to/backup/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DenormalizeOADPPrefix(tt.normalizedPrefix)
			if result != tt.expected {
				t.Errorf("DenormalizeOADPPrefix(%q) = %q, want %q", tt.normalizedPrefix, result, tt.expected)
			}
		})
	}
}

func TestGetOADPPrefix(t *testing.T) {
	prefix := GetOADPPrefix()
	expected := "oadp-vmdp/"
	if prefix != expected {
		t.Errorf("GetOADPPrefix() = %q, want %q", prefix, expected)
	}
}

func TestIsOADPPrefix(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		expected bool
	}{
		{
			name:     "valid OADP prefix",
			prefix:   "oadp-vmdp/",
			expected: true,
		},
		{
			name:     "valid OADP prefix with path",
			prefix:   "oadp-vmdp/abc/def",
			expected: true,
		},
		{
			name:     "invalid prefix",
			prefix:   "something-else/",
			expected: false,
		},
		{
			name:     "empty string",
			prefix:   "",
			expected: false,
		},
		{
			name:     "partial match",
			prefix:   "oadp-vm",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsOADPPrefix(tt.prefix)
			if result != tt.expected {
				t.Errorf("IsOADPPrefix(%q) = %v, want %v", tt.prefix, result, tt.expected)
			}
		})
	}
}

func TestRoundTripNormalizeDenormalize(t *testing.T) {
	tests := []string{
		"abc/",
		"abc",
		"path/to/backup/",
		"",
		"single",
	}

	for _, userPrefix := range tests {
		t.Run("roundtrip_"+userPrefix, func(t *testing.T) {
			normalized, err := NormalizeOADPPrefix(userPrefix)
			if err != nil {
				t.Fatalf("NormalizeOADPPrefix(%q) returned error: %v", userPrefix, err)
			}

			denormalized := DenormalizeOADPPrefix(normalized)
			if denormalized != userPrefix {
				t.Errorf("Round-trip failed: %q -> %q -> %q", userPrefix, normalized, denormalized)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}
