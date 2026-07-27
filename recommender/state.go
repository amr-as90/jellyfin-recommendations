// Package recommender contains the definitions and functions for the recommendation service
package recommender

import (
	"fmt"
	"sync"
)

type StateManager struct {
	Mu              sync.RWMutex
	Users           map[string]string
	UserFavorites   map[string]map[string]bool
	APIKey          string
	ServerURL       string
	Namer           *CollectionNamer
	UserCollections map[string]*UserCollectionInfo
}

// NewStateManager initializes a state manager
func NewStateManager(serverURL, apiKey string) *StateManager {
	return &StateManager{
		ServerURL:     serverURL,
		APIKey:        apiKey,
		UserFavorites: make(map[string]map[string]bool),
	}
}

// Sync performs initial server setup and state reconciliation
func (s *StateManager) Sync() error {
	// Get the language to name the new collections
	namer, err := s.NewCollectionNamerFromState()
	if err != nil {
		// Log warning and fallback to default rather than stopping startup
		namer = &CollectionNamer{Language: "en"}
	}
	s.Namer = namer

	// Get the users from Jellyfin
	users, err := s.getUsers()
	if err != nil {
		return fmt.Errorf("failed to get users from Jellyfin: %w", err)
	}

	// Hydrate user map
	s.Mu.Lock()
	for _, user := range users {
		s.Users[user.ID] = user.Name
	}
	s.Mu.Unlock()

	// Hydrate the map of user favorites
	err = s.hydrateFavorites(users)
	if err != nil {
		return fmt.Errorf("failed to hydrate user favorites: %w", err)
	}

	// Reconcile any differences
	err = s.reconcileCollections(users)
	if err != nil {
		return fmt.Errorf("failed to reconcile collections: %w", err)
	}

	return nil
}
