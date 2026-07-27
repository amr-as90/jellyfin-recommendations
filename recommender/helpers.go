package recommender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// GetJellyfin sends an authenticated GET request to the given Jellyfin endpoint
// and automatically parses the JSON response into type T.
func getJellyfin[T any](s *StateManager, endpoint string) (T, error) {
	var result T

	// Construct full URL
	reqURL := fmt.Sprintf("%s%s", s.ServerURL, endpoint)

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return result, fmt.Errorf("failed to create GET request for %s: %w", endpoint, err)
	}

	// Attach authentication header
	req.Header.Set("X-Emby-Token", s.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("HTTP request failed for %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("GET %s returned non-200 status: %s", endpoint, resp.Status)
	}

	// Decode JSON directly into result
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("failed to parse JSON response from %s: %w", endpoint, err)
	}

	return result, nil
}

// postJellyfin sends an authenticated POST request with an optional JSON body (payload)
// and unmarshals the JSON response into type T.
func postJellyfin[T any](s *StateManager, endpoint string, payload any) (T, error) {
	var result T

	reqURL := fmt.Sprintf("%s%s", s.ServerURL, endpoint)

	var bodyReader io.Reader
	if payload != nil {
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			return result, fmt.Errorf("failed to marshal POST payload for %s: %w", endpoint, err)
		}
		bodyReader = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(http.MethodPost, reqURL, bodyReader)
	if err != nil {
		return result, fmt.Errorf("failed to create POST request for %s: %w", endpoint, err)
	}

	req.Header.Set("X-Emby-Token", s.APIKey)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("POST HTTP request failed for %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return result, fmt.Errorf("POST %s returned status: %s", endpoint, resp.Status)
	}

	// If the API endpoint returns an empty body or zero content length, return result without decoding
	if resp.ContentLength == 0 || resp.StatusCode == http.StatusNoContent {
		return result, nil
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		// Ignore EOF if response body happens to be empty
		if err == io.EOF {
			return result, nil
		}
		return result, fmt.Errorf("failed to parse POST JSON response from %s: %w", endpoint, err)
	}

	return result, nil
}

// deleteJellyfin sends an authenticated DELETE request to the given endpoint.
func (s *StateManager) deleteJellyfin(endpoint string) error {
	reqURL := fmt.Sprintf("%s%s", s.ServerURL, endpoint)

	req, err := http.NewRequest(http.MethodDelete, reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request for %s: %w", endpoint, err)
	}

	req.Header.Set("X-Emby-Token", s.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("DELETE HTTP request failed for %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("DELETE %s returned status: %s", endpoint, resp.Status)
	}

	return nil
}
