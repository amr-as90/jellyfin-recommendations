// Package models contains all of the structs required
package models

import "sync"

type StateManager struct {
	Mu            sync.RWMutex
	UserFavorites map[string]map[string]bool
	PlaylistID    string
	APIKey        string
	ServerURL     string
}
