# Telescope

**Telescope** is a smart adaptive bitrate (ABR) proxy system designed for streaming content over the InterPlanetary File System (IPFS). It dynamically adjusts video quality based on network conditions and cache awareness, enhancing the efficiency and user experience of decentralized video delivery.

---

## 🚀 Features

- **Dynamic ABR Logic**: Rewrites DASH MPDs in real-time using IPFS, Gateway, and Client bandwidth estimations.
- **Cache Awareness**: Tracks segment caching status to optimize quality selection.
- **Stateless Microservices**: Replaces stateful proxy servers with scalable, stateless services.
- **Modern Observability**: Real-time metrics with Prometheus and distributed tracing with OpenTelemetry.
- **DASH.js Integration**: Fully compatible with DASH.js for live testing of streaming behavior.
- **Improved Architecture**: Clean, modular project structure with faster service using Go-Fiber.
- **Scalable Proxy System**: Designed to handle high traffic and large-scale deployments.

---

## 🆕 What’s New in This Version?

- Migration from **Gin** to **Go-Fiber** for improved performance.
- Redesigned ABR logic aligned with formal analysis and research papers.
- Enhanced observability with **Prometheus metrics** and **OpenTelemetry tracing**.
- File-based segment **cache tracking system** for better cache management.
- Support for **stateless microservices** to improve scalability and reliability.

---

## 📖 How Telescope Works

Telescope acts as a smart proxy between a DASH video player and IPFS. It intercepts video segment requests, estimates network conditions and caching status, and dynamically rewrites the video manifest (MPD) to guide adaptive quality selection.

### Workflow

1. **Manifest Request**:
   - The DASH.js player requests the `.mpd` manifest from Telescope.
   - Telescope fetches the segment CID list from IPFS and rewrites the MPD based on bandwidth estimations (`Tc`, `Tg`, `Tn`).
   - The rewritten MPD is returned to the player.

2. **Segment Request**:
   - The player requests video segments from Telescope.
   - Telescope fetches the segment from IPFS (or serves it from the cache if available).
   - The segment is streamed back to the player.

### Sequence Diagram

```plaintext
Client->>Proxy: GET /videos/bunny.mpd
Proxy->>IPFS: Fetch segment CID list
Proxy->>Proxy: Rewrites MPD based on Tc, Tg, Tn
Proxy-->>Client: Returns adaptive MPD

Client->>Proxy: GET /videos/bunny_128256bps/seg1.m4s
Proxy->>IPFS: Fetch segment or serve cached
Proxy-->>Client: Stream segment
```

![](.github/assets/diagram.svg)
![](.github/assets/sequence.svg)

---

## 📊 Metrics and Observability

Telescope provides real-time metrics and distributed tracing to support debugging, monitoring, and performance optimization.

### Key Metrics

- **RTT (Round-Trip Time)**: Measures segment fetch latency.
- **Bandwidth Estimation**: Calculates bandwidth based on segment size and transfer time.
- **Throughput Tracking**: Tracks client-reported throughput via HTTP headers.
- **Cache Awareness**: Monitors per-segment cache hit/miss ratios.
- **Segment Quality History**: Logs quality level shifts over time.

### Observability Tools

- **Prometheus**: Exposes metrics at `:9090/metrics`.
- **OpenTelemetry**: Provides distributed tracing with full support for Jaeger.

---

## 📂 Project Structure

```
proxy/
├── cmd/                     # Main entry points for the proxy
├── internal/                # Core application logic
│   ├── controllers/         # ABR logic and MPD rewriting
│   ├── storage/             # Cache management
│   └── metrics/             # Metrics and observability
services/                    # IPFS, Prometheus, Bootstrap, and Proxy config files
docker-compose.yaml          # Execute project using Docker
```

---

## 📈 Future Improvements

- **Multi-Replica Awareness**: Enhance ABR logic to account for multiple replicas in IPFS.
- **Advanced Caching Policies**: Implement predictive caching based on access patterns.
- **Support for HLS**: Extend support to HLS manifests in addition to DASH.
- **Improved Load Balancing**: Optimize proxy performance under high traffic.

---

## Run using Docker

1. First run the following scripts to download and encode videos.
   1. `scripts/videos/fetch.sh`
   2. `scripts/videos/encode.sh`
2. Then run `docker-compose up -d ipfs0 ipfs1 ipfs2`
3. After that run cluster bootstrap script `scripts/cluster/bootstrap.sh`
4. Then run `docker-compose up -d bootstrap`
5. After that you can run other services `docker-compose up -d telescope prometheus jaeger`

### UIs

- `localhost:5050` : Telescope proxy UI
- `localhost:9090` : Prometheus UI
- `localhost:16686` : Jaeger UI

---

## 📜 License

This project is licensed under the [MIT License](LICENSE).
