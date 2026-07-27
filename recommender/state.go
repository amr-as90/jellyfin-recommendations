// Package recommender contains the definitions and functions for the recommendation service
package recommender

import (
	"fmt"
	"sync"
)

type StateManager struct {
	Mu            sync.RWMutex
	UserFavorites map[string]map[string]bool
	APIKey        string
	ServerURL     string
	Namer         *CollectionNamer
}

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

	// Hydrate the map of user favorites
	err = s.hydrateFavorites()
	if err != nil {
		return fmt.Errorf("failed to hydrate user favorites: %w", err)
	}

	// Reconcile any differences

	return nil
}
