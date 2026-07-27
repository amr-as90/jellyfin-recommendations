package recommender

type UserCollectionInfo struct {
	CollectionID string
	ItemIDs      map[string]bool
}

// User represents a Jellyfin user account
type user struct {
	ID   string `json:"Id"`
	Name string `json:"Name"`
}

// BaseItem represents an item in the library (movie, show, etc.,)
type baseItem struct {
	ID   string `json:"Id"`
	Name string `json:"Name"`
}

// QueryUserFavoritesResponse represents the payload returned by GET /Users/{userId}/Items
type queryUserFavoritesResponse struct {
	Items []baseItem `json:"Items"`
}

// SystemConfiguration is just the language portion of the server's system configuration so that
// we know how to name the recommendation collections
type systemConfiguration struct {
	UICulture                 string `json:"UICulture"`
	PreferredMetadataLanguage string `json:"PreferredMetadataLanguage"`
}

// ExistingCollectionsResponse maps the JSON response from GET /Items?includeItemTypes=BoxSet
type existingCollectionsResponse struct {
	Items []struct {
		ID   string `json:"Id"`
		Name string `json:"Name"`
	} `json:"Items"`
}

// createCollectionResponse captures the payload returned when creating a collection
type createCollectionResponse struct {
	ID string `json:"Id"`
}
