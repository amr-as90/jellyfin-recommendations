package services

// BaseItem represents an item returned from Jellyfin query responses
type BaseItem struct {
	ID   string `json:"Id"`
	Name string `json:"Name"`
}

// QueryItemsResponse represents the wrapper returned by GET /Items
type QueryItemsResponse struct {
	Items []BaseItem `json:"Items"`
}

// CreatePlaylistRequest represents the payload expected by POST /Playlists
type CreatePlaylistRequest struct {
	Name       string   `json:"Name"`
	OpenAccess bool     `json:"OpenAccess"`
	ItemIDs    []string `json:"Ids"`
}

// CreatePlaylistResponse represents the object returned when a playlist is created
type CreatePlaylistResponse struct {
	ID string `json:"Id"`
}
