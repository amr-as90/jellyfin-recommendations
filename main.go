package main

import (
	"log"
	"os"

	"github.com/amr-as90/jellyfin-recommendations/helpers"
	"github.com/amr-as90/jellyfin-recommendations/models"
)

func main() {
	// Load the API Key from the environment variable
	serverURL := os.Getenv("JELLYFIN_URL")
	apiKey := os.Getenv("API_KEY")
	playlistID := os.Getenv("PLAYLIST_ID")

	if serverURL == "" || apiKey == "" {
		log.Fatal("A valid Jellyfin URL and API Key are required")
	}

	normalizedURL, err := helpers.ValidateServerURL(serverURL)
	if err != nil {
		log.Fatal("Server URL isn't valid: %w", err)
	}

	state := &models.StateManager{
		APIKey:     apiKey,
		ServerURL:  normalizedURL,
		PlaylistID: playlistID,
	}

	// Check if we have a recorded playlist ID (Possibly just get this from the ENV variable too)

	// If we don't have a playlist ID, we create a playlist and store it's ID in the playlist ID variable

	// Get the user favorites and store them in a map

	// Check the playlist for user favorites, resolve any discrepancies (both favorites that don't exist, and items that exist which are not favorites)

	// Start the WebSocket connection and listen for favorite changes

	// If a user favorites or unfavorites, we check the map, if the action was a favorite and the item doesn't already exist, we add it. If the action
	// is to unfavorite, we check to see if anyone else has it favorited, if not, then we remove it from the playlist.
}
