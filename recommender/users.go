package recommender

import (
	"fmt"
	"log"
)

func (s *StateManager) hydrateFavorites(users []user) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	// Get user favorites
	for _, user := range users {
		favorites, err := s.getUserFavorites(user.ID)
		if err != nil {
			log.Printf("Warning: failed to fetch favorites for user %s (%s): %v", user.Name, user.ID, err)
			continue
		}

		// Map favorited items to this user
		for _, item := range favorites {
			if _, exists := s.UserFavorites[item.ID]; !exists {
				s.UserFavorites[item.ID] = make(map[string]bool)
			}
			s.UserFavorites[item.ID][user.ID] = true
		}
	}

	return nil
}

func (s *StateManager) getUsers() (users []user, err error) {
	return getJellyfin[[]user](s, "/Users")
}

func (s *StateManager) getUserFavorites(userID string) ([]baseItem, error) {
	if userID == "" {
		return nil, fmt.Errorf("received empty userID, cannot be empty")
	}

	endpoint := fmt.Sprintf("%s/Users/%s/Items?filters=IsFavorite&recursive=true", s.ServerURL, userID)

	resp, err := getJellyfin[queryUserFavoritesResponse](s, endpoint)
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}
