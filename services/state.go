// Package services contains the definitions and functions for the recommendation service
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

type StateManager struct {
	Mu            sync.RWMutex
	UserFavorites map[string]map[string]bool
	PlaylistID    string
	APIKey        string
	ServerURL     string
}

// GetOrCreatePlaylist queries Jellyfin for a playlist with our name, and returns its ID if it exists,
// otherwise, it creates a new playlist with our name and returns its ID.
func (s *StateManager) GetOrCreatePlaylist(playlistName string) (err error) {
	// Query existing playlists
	playlists, err := s.getPlaylists()
	if err != nil {
		return fmt.Errorf("couldn't get existing playlists: %w", err)
	}

	// Check if a playlist with our playlist name exists
	for _, item := range playlists.Items {
		if item.Name == playlistName {
			s.PlaylistID = item.ID
			return nil
		}
	}

	// Otherwise, we create one
	err = s.createPlaylist(playlistName)
	if err != nil {
		return fmt.Errorf("couldn't create new playlist: %w", err)
	}

	return nil
}

func (s *StateManager) getPlaylists() (playlists QueryPlaylistsResponse, err error) {
	reqURL := fmt.Sprintf("%s/Items?includeItemTypes=Playlist&recursive=true", s.ServerURL)

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return playlists, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Emby-Token", s.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return playlists, fmt.Errorf("failed to query playlists: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return playlists, fmt.Errorf("querying playlists returned a non-200 status: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&playlists); err != nil {
		return playlists, fmt.Errorf("failed to parse playlist query response: %w", err)
	}

	return playlists, nil
}

func (s *StateManager) createPlaylist(playlistName string) (err error) {
	createURL := fmt.Sprintf("%s/Playlists?name=%s&isOpenAccess=true", s.ServerURL, url.QueryEscape(playlistName))

	createReqPayload := CreatePlaylistRequest{
		Name:       playlistName,
		OpenAccess: true,
		ItemIDs:    []string{},
	}

	bodyBytes, err := json.Marshal(createReqPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal create playlist request: %w", err)
	}

	createReq, err := http.NewRequest(http.MethodPost, createURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create playlist request: %w", err)
	}

	createReq.Header.Set("X-Emby-Token", s.APIKey)
	createReq.Header.Set("Content-Type", "application/json")

	createResp, err := http.DefaultClient.Do(createReq)
	if err != nil {
		return fmt.Errorf("error occurred when attempting to create playlist: %w", err)
	}
	defer createResp.Body.Close()

	if createResp.StatusCode != http.StatusOK && createResp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(createResp.Body)
		return fmt.Errorf("failed to create playlist (status %s): %s", createResp.Status, string(respBody))
	}

	var createResult CreatePlaylistResponse
	if err := json.NewDecoder(createResp.Body).Decode(&createResult); err != nil {
		return fmt.Errorf("failed to parse create playlist response: %w", err)
	}

	// Store the newly created playlist ID on our state manager
	s.PlaylistID = createResult.ID
	return nil
}
