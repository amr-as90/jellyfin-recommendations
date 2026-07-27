# Jellyfin Recommendations

Jellyfin recommendations is a simple, lightweight program written in Go, which turns user favorites into a 'recommendations' collection.

![Jellyfin screenshot showing recommendations collection](/images/jellyfin1.avif)

The collection is divided into user names. If a user has some items as favorites, they show up under their username as recommendations within the collection.

![Jellyfin screenshot showing the different user's favorites](/images/jellyfin2.avif)

Inside, it simply uses the already existing collections functionality to show you their recommended movies and TV shows.

![Jellyfin screenshot showing recommendations as user favorites](/images/jellyfin3.avif)

*Note: It runs every 20 minutes. Initially I really wanted to get this working with WebSockets, but it doesn't seem like this information is passed via WebSocket, and its too much of a pain on the user side to get this working with WebHooks, so I figured since real-time updates aren't necessary, it just syncs at an interval.*

# Deployment

```bash
docker run -d \
  --name jellyfin-recommendations \
  --restart unless-stopped \
  -e JELLYFIN_URL="{YOUR_JELLYFIN_URL}" \
  -e JELLYFIN_API_KEY="{YOUR_API_KEY}" \
  ghcr.io/amr-as90/jellyfin-recommendations:latest
```

Replace `{YOUR_JELLYFIN_URL}` and `{YOUR_API_KEY}` with your actual Jellyfin URL and API key.