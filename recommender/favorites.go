package recommender

import "fmt"

func (s *StateManager) addUserFavorite(userID, userName, itemID string) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	// Update in-memory favorites
	if s.UserFavorites[itemID] == nil {
		s.UserFavorites[itemID] = make(map[string]bool)
	}
	s.UserFavorites[itemID][userID] = true

	col, exists := s.UserCollections[userID]

	// If no collection exists, create it
	if !exists || len(col.ItemIDs) == 0 {
		collectionName := userName
		colID, err := s.createCollection(collectionName, itemID)
		if err != nil {
			return err
		}

		s.UserCollections[userID] = &UserCollectionInfo{
			CollectionID: colID,
			ItemIDs:      map[string]bool{itemID: true},
		}
		return nil
	}

	// Collection already exists, add the item
	if !col.ItemIDs[itemID] {
		if err := s.addItemToCollection(col.CollectionID, itemID); err != nil {
			return err
		}
		col.ItemIDs[itemID] = true
	}

	return nil
}

func (s *StateManager) removeUserFavorite(userID, itemID string) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	// Remove from in-memory favorites
	if s.UserFavorites[itemID] != nil {
		delete(s.UserFavorites[itemID], userID)
	}

	col, exists := s.UserCollections[userID]
	if !exists || !col.ItemIDs[itemID] {
		return nil
	}

	// Remove item from Jellyfin collection
	if err := s.removeItemFromCollection(col.CollectionID, itemID); err != nil {
		return err
	}
	delete(col.ItemIDs, itemID)

	// If collection is now empty, delete it
	if len(col.ItemIDs) == 0 {
		endpoint := fmt.Sprintf("/Items/%s", col.CollectionID)
		if err := s.deleteJellyfin(endpoint); err != nil {
			return fmt.Errorf("failed to delete now-empty collection: %w", err)
		}
		delete(s.UserCollections, userID)
	}

	return nil
}
