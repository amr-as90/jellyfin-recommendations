package recommender

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func (s *StateManager) reconcileCollections(users []user) error {
	existingCollections, err := s.getExistingCollections()
	if err != nil {
		return fmt.Errorf("reconciliation failed to query existing collections: %w", err)
	}

	s.Mu.Lock()
	defer s.Mu.Unlock()

	// Reset collection map
	s.UserCollections = make(map[string]*UserCollectionInfo)

	for _, user := range users {
		expectedName := user.Name
		collectionID, exists := existingCollections[expectedName]

		// Get all expected item IDs for this user based on their favorites
		expectedItems := make(map[string]bool)
		for itemID, userMap := range s.UserFavorites {
			if userMap[user.ID] {
				expectedItems[itemID] = true
			}
		}

		// If the user has favorites
		if len(expectedItems) > 0 {
			var expectedSlice []string
			for id := range expectedItems {
				expectedSlice = append(expectedSlice, id)
			}

			// Create collection if missing
			if !exists {
				firstItem := expectedSlice[0]
				newID, err := s.createCollectionWithImage(expectedName, user.ID, firstItem)
				if err != nil {
					log.Printf("Failed to create collection %q for user %s: %v", expectedName, user.Name, err)
					continue
				}
				collectionID = newID
				expectedSlice = expectedSlice[1:]
			}

			// Get the items in the collection from Jellyfin
			actualItems, err := s.getCollectionItemIDs(collectionID)
			if err != nil {
				log.Printf("Failed to read items for collection %s: %v", expectedName, err)
				actualItems = make(map[string]bool)
			}

			// Add missing items
			for _, itemID := range expectedSlice {
				if !actualItems[itemID] {
					if err := s.addItemToCollection(collectionID, itemID); err != nil {
						log.Printf("Failed to add item %s: %v", itemID, err)
					} else {
						actualItems[itemID] = true
					}
				}
			}

			// Remove stale items
			for actualID := range actualItems {
				if !expectedItems[actualID] {
					if err := s.removeItemFromCollection(collectionID, actualID); err != nil {
						log.Printf("Failed to remove item %s: %v", actualID, err)
					} else {
						delete(actualItems, actualID)
					}
				}
			}

			// Save cached state for runtime
			s.UserCollections[user.ID] = &UserCollectionInfo{
				CollectionID: collectionID,
				ItemIDs:      actualItems,
			}

		} else {
			// Otherwise, if the user doesn't have favorites
			if exists {
				endpoint := fmt.Sprintf("/Items/%s", collectionID)
				if err := s.deleteJellyfin(endpoint); err != nil {
					log.Printf("Failed to delete empty collection %s: %v", expectedName, err)
				}
			}
		}
	}

	return nil
}

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

func (s *StateManager) createCollectionWithImage(name, userID, initialItemID string) (string, error) {
	newID, err := s.createCollection(name, initialItemID)
	if err != nil {
		return "", err
	}

	imageBytes, contentType, err := s.getUserProfilePicture(userID)
	if err != nil {
		return "", err
	}

	err = s.setCollectionImage(newID, imageBytes, contentType)
	if err != nil {
		return "", err
	}

	return newID, nil
}

func (s *StateManager) createCollection(name, initialItemID string) (string, error) {
	if name == "" || initialItemID == "" {
		return "", fmt.Errorf("collection name and initialItemID cannot be empty")
	}

	params := url.Values{}
	params.Set("Name", name)
	params.Set("Ids", initialItemID)

	endpoint := fmt.Sprintf("/Collections?%s", params.Encode())

	resp, err := postJellyfin[createCollectionResponse](s, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create collection %q: %w", name, err)
	}

	return resp.ID, nil
}

func (s *StateManager) addItemToCollection(collectionID, itemID string) error {
	if collectionID == "" || itemID == "" {
		return fmt.Errorf("collectionID and itemID cannot be empty")
	}

	params := url.Values{}
	params.Set("Ids", itemID)

	endpoint := fmt.Sprintf("/Collections/%s/Items?%s", collectionID, params.Encode())

	_, err := postJellyfin[any](s, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to add item %s to collection %s: %w", itemID, collectionID, err)
	}

	return nil
}

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

func (s *StateManager) getUserProfilePicture(userID string) ([]byte, string, error) {
	profilePictureURL := fmt.Sprintf("%s/Users/%s/Images/Primary", s.ServerURL, userID)

	req, err := http.NewRequest(http.MethodGet, profilePictureURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("X-Emby-Token", s.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("avatar fetch returned HTTP %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return data, contentType, nil
}

func (s *StateManager) setCollectionImage(itemID string, imageData []byte, contentType string) error {
	endpoint := fmt.Sprintf("%s/Items/%s/Images/Primary", s.ServerURL, itemID)

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(imageData))
	if err != nil {
		return fmt.Errorf("failed to create POST image request: %w", err)
	}

	req.Header.Set("X-Emby-Token", s.APIKey)
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("POST image request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("POST image returned status: %s", resp.Status)
	}

	return nil
}
