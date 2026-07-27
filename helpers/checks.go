// Package helpers contains helper functions
package helpers

import (
	"fmt"
	"net/url"
	"strings"
)

// ValidateServerURL accepts a URL and returns a normalized URL or an error if it is invalid
func ValidateServerURL(rawURL string) (normalizedURL string, err error) {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return "", fmt.Errorf("server URL is empty")
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "", fmt.Errorf("invalid URL format: %w", err)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("invalid URL scheme %q: URL must start with http:// or https://", parsed.Scheme)
	}

	if parsed.Host == "" {
		return "", fmt.Errorf("URL must include a host (e.g., localhost:8096 or 192.168.1.50)")
	}

	normalized := strings.TrimRight(trimmed, "/")

	return normalized, nil
}
