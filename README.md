# Telescope

**Telescope** is a proxy service that works like a smart adaptive bitrate (ABR) system for the InterPlanetary File System (IPFS). It helps manage and improve how content is delivered from IPFS by adjusting the quality based on network conditions.

This project is a rebuilt version of the original Telescope, which was cloned from [github.com/SBUNetSys/Telescope](https://github.com/SBUNetSys/Telescope). Our version includes several improvements:

- Cleaner and more efficient source code, written in **Golang**  
- A better **monitoring system** to track performance and activity  
- An improved **caching system** to make content delivery faster and more reliable  

Telescope is designed to make IPFS-based streaming and content delivery smarter, smoother, and more responsive to different network speeds.

## How Video Streaming Works with Telescope?

Telescope acts as a middle layer between the video player and the IPFS network. It takes requests from the player, fetches video segments from IPFS, and decides which video quality to serve based on the current network speed (just like adaptive bitrate streaming).

The process looks like this:

1. **The user plays a video** in their browser or media player.
2. **The player sends a request** for video segments to Telescope.
3. **Telescope checks the network conditions** and chooses the best quality (bitrate) for smooth playback.
4. **It fetches the correct segment** from IPFS.
5. **The segment is sent back** to the player for viewing.
6. This cycle repeats for each segment of the video, adjusting the quality if the network speed changes.

This approach ensures a better viewing experience with less buffering and faster loading times, even when the network is unstable.

```
Client->>Proxy: GET /dash/bafy123
Proxy->>IPFS: Fetch metadata from bafy123
IPFS-->>Proxy: Returns video metadata
Proxy->>Proxy: Generate MPD with embedded CID
Proxy-->>Client: Returns DASH manifest

Client->>Proxy: GET /segment/bafy123/0-999999
Proxy->>IPFS: cat bafy123 (range 0-999999)
IPFS-->>Proxy: Returns segment data
Proxy-->>Client: Serves video segment
```

### Usage Diagram

![](.github/assets/diagram.svg)

### Sequence Diagram

![](.github/assets/sequence.svg)

## Metrics

Telescope uses several key metrics to make smart decisions about video quality. These help it choose the best bitrate for smooth and efficient streaming.

Here’s what it takes into account:

- **RTT (Round-Trip Time):**  
  Telescope measures the time it takes to send a request to an IPFS node and get a response. This is done for every segment request.

- **Bandwidth Estimation:**  
  Bandwidth is calculated based on the size of the video segment and the RTT. This helps estimate how fast the network is at any moment.

- **Segment Count:**  
  The number of segments requested and delivered is also tracked to help make more accurate streaming decisions.

- **Caching Status:**  
  Telescope checks if the requested file is already cached. If it is, it can deliver it much faster, which improves performance.

By using these metrics, Telescope adapts the video quality in real-time to match current network conditions, giving users a smoother viewing experience.
