package recommender

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// StartWebSocketListener connects to Jellyfin's WebSocket and streams favorite updates.
func (s *StateManager) StartWebSocketListener() error {
	// Convert http:// or https:// to ws:// or wss://
	wsURL := strings.Replace(s.ServerURL, "http://", "ws://", 1)
	wsURL = strings.Replace(wsURL, "https://", "wss://", 1)

	fullURL := fmt.Sprintf("%s/socket?api_key=%s", wsURL, s.APIKey)

	conn, _, err := websocket.DefaultDialer.Dial(fullURL, http.Header{})
	if err != nil {
		return fmt.Errorf("websocket connection failed: %w", err)
	}
	defer conn.Close()

	// TODO: Replace these with Slog messages
	log.Println("Successfully connected to Jellyfin WebSocket.")

	// Keep-alive / initial subscription message
	subMsg := `{"MessageType":"SessionsStart", "Data":"0,1000"}`
	if err := conn.WriteMessage(websocket.TextMessage, []byte(subMsg)); err != nil {
		return fmt.Errorf("failed to send subscription message: %w", err)
	}

	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v. Reconnecting...", err)
			return err
		}

		var msg WSMessage
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			continue
		}

		if msg.MessageType == "UserDataChanged" {
			s.handleUserDataChanged(msg.Data)
		}
	}
}

func (s *StateManager) handleUserDataChanged(data UserDataChanged) {
	for _, change := range data.UserDataList {
		userID := change.UserID
		itemID := change.ItemID
		isFav := change.Data.IsFavorite

		// Find user name for collection formatting (from a quick cache or user list)
		userName := s.getUserNameByID(userID)
		if userName == "" {
			continue
		}

		if isFav {
			log.Printf("User %s favorited item %s", userName, itemID)
			if err := s.addUserFavorite(userID, userName, itemID); err != nil {
				log.Printf("Error adding favorite for %s: %v", userName, err)
			}
		} else {
			log.Printf("User %s UNfavorited item %s", userName, itemID)
			if err := s.removeUserFavorite(userID, itemID); err != nil {
				log.Printf("Error removing favorite for %s: %v", userName, err)
			}
		}
	}
}
