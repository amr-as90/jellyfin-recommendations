// Package services contains the definitions and functions for the recommendation service
package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type StateManager struct {
	Mu            sync.RWMutex
	UserFavorites map[string]map[string]bool
	PlaylistID    string
	APIKey        string
	ServerURL     string
}

func (s *StateManager) HydrateFavorites() error {
	// Get all users
	users, err := s.getUsers()
	if err != nil {
		return err
	}

	s.Mu.Lock()
	defer s.Mu.Unlock()

	// Get user favorites
	for _, user := range users {
		favorites, err := s.getUserFavorites(user.ID)
		if err != nil {
			log.Printf("Warning: failed to fetch favorites for user %s (%s): %v", user.Name, user.ID, err)
			continue
		}

		// Map favorited items to this user
		for _, item := range favorites {
			if _, exists := s.UserFavorites[item.ID]; !exists {
				s.UserFavorites[item.ID] = make(map[string]bool)
			}
			s.UserFavorites[item.ID][user.ID] = true
		}
	}

	return nil
}

func (s *StateManager) getUsers() (users []User, err error) {
	reqURL := fmt.Sprintf("%s/Users", s.ServerURL)

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a request to get users: %w", err)
	}
	req.Header.Set("X-Emby-Token", s.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return users, fmt.Errorf("failed to execute request to get users: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get users returned non-200 status: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		return nil, fmt.Errorf("failed to parse users response: %w", err)
	}

	return users, nil
}

func (s *StateManager) getUserFavorites(userID string) ([]BaseItem, error) {
	if userID == "" {
		return nil, fmt.Errorf("received empty userID, cannot be empty")
	}

	reqURL := fmt.Sprintf("%s/Users/%s/Items?filters=IsFavorite&recursive=true", s.ServerURL, userID)

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user favorites request: %w", err)
	}

	req.Header.Set("X-Emby-Token", s.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute user favorites request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get user favorites returned non-200 status: %s", resp.Status)
	}

	var favResponse QueryUserFavoritesResponse
	if err := json.NewDecoder(resp.Body).Decode(&favResponse); err != nil {
		return nil, fmt.Errorf("failed to parse user favorites response: %w", err)
	}

	return favResponse.Items, nil
}
