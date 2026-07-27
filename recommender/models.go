package recommender

// User represents a Jellyfin user account
type User struct {
	ID   string `json:"Id"`
	Name string `json:"Name"`
}

// BaseItem represents an item in the library (movie, show, etc.,)
type BaseItem struct {
	ID   string `json:"Id"`
	Name string `json:"Name"`
}

// QueryUserFavoritesResponse represents the payload returned by GET /Users/{userId}/Items
type QueryUserFavoritesResponse struct {
	Items []BaseItem `json:"Items"`
}

type SystemConfiguration struct {
	UICulture                 string `json:"UICulture"`
	PreferredMetadataLanguage string `json:"PreferredMetadataLanguage"`
}
