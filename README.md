# Jellyfin Recommendations

**Jellyfin Recommendations** is a simple, lightweight program which hijacks the 'Collections' feature in Jellyfin to turn user favorites into server recommendations.

![Jellyfin screenshot showing recommendations collection](/images/jellyfin1.avif)

Each user on the server (that has favorites) appears as a 'Collection'. Each collection contains that user's favorites.

![Jellyfin screenshot showing the different user's favorites](/images/jellyfin2.avif)

Inside, we simply use the already existing 'Collections' functionality to neatly show you that user's recommended movies and TV shows.

![Jellyfin screenshot showing recommendations as user favorites](/images/jellyfin3.avif)

This has the added benefit of working on all native clients, since we are 'taking over' an existing feature, and not injecting any custom JavaScript, etc.,

*Note: It runs every 20 minutes. Initially I wanted to get this working with WebSockets, but it doesn't seem like user-favorite information is passed via WebSocket, and for me, it feels like too much friction on the user's side to get this working with WebHooks. I figured since real-time updates aren't necessary, it just syncs at an interval.*

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

## Optional Step

There is only one native 'Collections' tab in Jellyfin. If you don't really use collections like I do, you can turn off the 'Group movies and shows into collections' setting in the Jellyfin dashboard, and when **Jellyfin Recommendations** populates the collections, simply rename it to 'Recommendations', or whatever you prefer, as this can't be done via the API as far as I can tell.