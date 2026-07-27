package main

import (
	"log"
	"os"
	"time"

	"github.com/amr-as90/jellyfin-recommendations/recommender"
)

func main() {
	// Load the environment variables
	serverURL := os.Getenv("JELLYFIN_URL")
	apiKey := os.Getenv("API_KEY")

	if serverURL == "" || apiKey == "" {
		log.Fatal("A valid Jellyfin URL and API Key are required")
	}

	normalizedURL, err := recommender.ValidateServerURL(serverURL)
	if err != nil {
		log.Fatal("Server URL isn't valid: %w", err)
	}

	state := recommender.NewStateManager(normalizedURL, apiKey)

	err = state.Sync()
	if err != nil {
		log.Fatal("Something went wrong during initial sync: %w", err)
	}

	for {
		err = state.StartWebSocketListener()
		if err != nil {
			log.Printf("WebSocket disconnected: %v. Retrying in 3 seconds...", err)
			time.Sleep(3 * time.Second)
		}
	}

}
