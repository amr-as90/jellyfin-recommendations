package recommender

import (
	"fmt"
)

func (s *StateManager) ReconcileCollections() error {
	return nil
}

// getCollectionItemIDs returns a set of item IDs inside a collection.
func (s *StateManager) getCollectionItemIDs(collectionID string) (map[string]bool, error) {
	endpoint := fmt.Sprintf("/Items?parentId=%s", collectionID)

	resp, err := getJellyfin[queryUserFavoritesResponse](s, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch collection items for %s: %w", collectionID, err)
	}

	itemSet := make(map[string]bool)
	for _, item := range resp.Items {
		itemSet[item.ID] = true
	}

	return itemSet, nil
}

// getExistingCollections queries Jellyfin for all existing collections
func (s *StateManager) getExistingCollections() (map[string]string, error) {
	endpoint := "/Items?includeItemTypes=BoxSet&recursive=true"

	resp, err := getJellyfin[existingCollectionsResponse](s, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch existing collections: %w", err)
	}

	collections := make(map[string]string)
	for _, item := range resp.Items {
		collections[item.Name] = item.ID
	}

	return collections, nil
}

// createCollection creates a new collection in Jellyfin initialized with a single item ID.
// Returns the newly created collection ID.
func (s *StateManager) createCollection(name, initialItemID string) (string, error) {
	if name == "" || initialItemID == "" {
		return "", fmt.Errorf("collection name and initialItemID cannot be empty")
	}

	endpoint := fmt.Sprintf("/Collections?Name=%s&Ids=%s", name, initialItemID)

	resp, err := postJellyfin[createCollectionResponse](s, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create collection %q: %w", name, err)
	}

	return resp.ID, nil
}

// removeItemFromCollection removes an item from an existing Jellyfin collection.
func (s *StateManager) removeItemFromCollection(collectionID, itemID string) error {
	if collectionID == "" || itemID == "" {
		return fmt.Errorf("collectionID and itemID cannot be empty")
	}

	endpoint := fmt.Sprintf("/Collections/%s/Items?Ids=%s", collectionID, itemID)

	if err := s.deleteJellyfin(endpoint); err != nil {
		return fmt.Errorf("failed to remove item %s from collection %s: %w", itemID, collectionID, err)
	}

	return nil
}
