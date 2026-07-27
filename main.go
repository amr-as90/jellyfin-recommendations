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
	log.Printf("Initial sync complete")

	ticker := time.NewTicker(time.Minute * 20)
	defer ticker.Stop()

	for range ticker.C {
		err = state.Sync()
		if err != nil {
			log.Printf("Something went wrong during sync: %v", err)
		} else {
			log.Printf("Sync complete")
		}
	}
}
