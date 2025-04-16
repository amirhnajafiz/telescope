# Telescope

**Telescope** is a smart adaptive bitrate (ABR) proxy system for streaming content over the InterPlanetary File System (IPFS). It dynamically adjusts video quality based on network conditions and cache awareness, improving the efficiency and experience of decentralized video delivery.

This project is a **rebuilt version** of the original [Telescope proxy](https://github.com/SBUNetSys/Telescope), redesigned from scratch with a clean architecture, modern observability tools, and improved ABR logic based on formal analysis.

---

##  What’s New in This Version?

- ✅ Modular and idiomatic **Golang project structure** (`cmd/`, `internal/`)
- ✅ Rewritten **throughput estimation engine** using exponential smoothing for cached and uncached segments
- ✅ Paper-aligned **ABR logic** using accurate bandwidth deltas (`Tc - Tn`, `Tc - Tg`) to rewrite DASH MPDs dynamically
- ✅ Real-time **Prometheus metrics** and **OpenTelemetry tracing**
- ✅ Full-featured **DASH.js browser client** to test streaming behavior and adaptation live
- ✅ In-memory segment **cache tracking system**
- ✅ Fully tested core modules with unit tests for ABR, cache, and estimator

---

## How Video Streaming Works with Telescope

Telescope acts as a smart proxy between a DASH video player and IPFS. It receives video segment requests, estimates network conditions and caching status, and rewrites the video manifest (MPD) in real-time to guide quality selection.

The flow:

1. A user opens a video player (DASH.js) in their browser.
2. The player requests the `.mpd` manifest from Telescope.
3. Telescope dynamically rewrites the MPD based on current client throughput and cache status.
4. As the player requests segments, Telescope fetches them from IPFS (or a stub in test mode), tracks bandwidth, and updates throughput estimation.
5. The next manifest request reflects this updated bandwidth info, helping DASH.js adapt quality accordingly.

```
Client->>Proxy: GET /videos/bunny.mpd
Proxy->>IPFS: Fetch segment CID list
Proxy->>Proxy: Rewrites MPD based on Tc, Tg, Tn
Proxy-->>Client: Returns adaptive MPD

Client->>Proxy: GET /videos/bunny_128256bps/seg1.m4s
Proxy->>IPFS: Fetch segment or serve cached
Proxy-->>Client: Stream segment
```

---

### Usage Diagram

![](.github/assets/diagram.svg)

### Sequence Diagram

![](.github/assets/sequence.svg)

---

## Metrics and Observability

Telescope exposes real-time metrics and traces to support debugging, monitoring, and performance tuning.

### Key Metrics

- **RTT (Round-Trip Time)** – segment fetch latency
- **Bandwidth Estimation** – based on segment size and transfer time
- **Throughput Tracking** – smoothed client `Tc`, cached `Tg`, uncached `Tn`
- **Cache Awareness** – per-segment cache hit/miss ratio
- **Active Connections** – current number of live users
- **Segment Quality History** – logs how quality levels shift over time

### Observability Tools

- **Prometheus**: exports metrics via `/metrics`
- **Jaeger/OTel**: full support for distributed tracing via OpenTelemetry

---

## Project Structure

```
telescope/
├── cmd/               → API registerer-AKA Serviec Injector (api.go)
├── internal/          → core modules (api, cache, abr, throughput, telemetry)
├── client/            → DASH.js HTML test client
├── assets/            → .mpd + .m4s video files (for local testing)
├── scripts/           → utils (downloaders, etc.)
```

---

