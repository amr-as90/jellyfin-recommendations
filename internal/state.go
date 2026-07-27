// Package recommender contains the definitions and functions for the recommendation service
package recommender

import (
	"sync"
)

type StateManager struct {
	Mu            sync.RWMutex
	UserFavorites map[string]map[string]bool
	PlaylistID    string
	APIKey        string
	ServerURL     string
}
