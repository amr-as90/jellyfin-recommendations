package main

import (
	"log"
	"os"

	"github.com/amr-as90/jellyfin-recommendations/recommender"
)

func main() {
	// Load the API Key from the environment variable
	serverURL := os.Getenv("JELLYFIN_URL")
	apiKey := os.Getenv("API_KEY")

	if serverURL == "" || apiKey == "" {
		log.Fatal("A valid Jellyfin URL and API Key are required")
	}

	normalizedURL, err := recommender.ValidateServerURL(serverURL)
	if err != nil {
		log.Fatal("Server URL isn't valid: %w", err)
	}

	state := &recommender.StateManager{
		APIKey:    apiKey,
		ServerURL: normalizedURL,
	}

	err = state.Sync()
	if err != nil {
		log.Fatal("Something went wrong hydrating favorites: %w", err)
	}

	// Start a WebSocket listener to handle new user favorites and items removed from user favorites
}
