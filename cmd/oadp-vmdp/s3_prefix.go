package main

import (
	"strings"

	"github.com/pkg/errors"
)

const (
	// oadpS3Prefix is the required prefix for all OADP S3 repositories.
	// This ensures isolation from other Kopia repositories in shared S3 buckets.
	oadpS3Prefix = "oadp-vmdp/"
)

// NormalizeOADPPrefix prepends "oadp-vmdp/" to the user-provided prefix.
// It handles various user input formats:
// - Errors out if user prefix contains "oadp-vmdp" (case-insensitive)
// - Removes leading slashes from user prefix
// - Always prepends "oadp-vmdp/" to the result
// - Preserves trailing slashes from user input
//
// Examples:
//   - "abc/" → "oadp-vmdp/abc/"
//   - "abc" → "oadp-vmdp/abc"
//   - "/abc/bce/" → "oadp-vmdp/abc/bce/"
//   - "" → "oadp-vmdp/"
//   - "oadp-vmdp/abc" → error
func NormalizeOADPPrefix(userPrefix string) (string, error) {
	// Check if user is trying to provide the oadp-vmdp prefix themselves
	lowerPrefix := strings.ToLower(userPrefix)
	if strings.Contains(lowerPrefix, "oadp-vmdp") {
		return "", errors.New("prefix must not contain 'oadp-vmdp' - this prefix is automatically added")
	}

	// Remove leading slashes from user prefix
	cleanedPrefix := strings.TrimLeft(userPrefix, "/")

	// Combine oadp prefix with cleaned user prefix
	return oadpS3Prefix + cleanedPrefix, nil
}

// DenormalizeOADPPrefix removes the "oadp-vmdp/" prefix from a normalized prefix.
// This is used when displaying prefixes to users or saving to configuration files.
//
// Examples:
//   - "oadp-vmdp/abc/" → "abc/"
//   - "oadp-vmdp/abc" → "abc"
//   - "oadp-vmdp/" → ""
//   - "something-else" → "something-else" (unchanged if no OADP prefix)
func DenormalizeOADPPrefix(normalizedPrefix string) string {
	return strings.TrimPrefix(normalizedPrefix, oadpS3Prefix)
}

// GetOADPPrefix returns the OADP S3 prefix constant.
// This is useful for testing and documentation purposes.
func GetOADPPrefix() string {
	return oadpS3Prefix
}

// IsOADPPrefix checks if a given prefix starts with the OADP prefix.
func IsOADPPrefix(prefix string) bool {
	return strings.HasPrefix(prefix, oadpS3Prefix)
}
