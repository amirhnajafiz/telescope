# Telescope Proxy

The purpose of this proxy is to serve as a gateway for uploading and streaming videos over IPFS. It is a web server implemented using the Golang Fiber framework. Since the proxy server is stateless, it can be easily scaled to accommodate varying workloads. Caching is handled through HTTP headers on the client side, ensuring efficient data retrieval.  

This gateway proxy enhances the efficiency and performance of Adaptive Bitrate (ABR) streaming protocols. Additionally, as a gateway, it provides performance metrics and tracing capabilities, facilitating better monitoring and observability.  

The original Telescope proxy is available at [Telescope Repository](https://github.com/SBUNetSys/Telescope). This project is a rebuilt version of Telescope, designed to address its existing issues while introducing new features and improvements.

---

## Improvements & Progress

+ Re-architected project following idiomatic Go structure (`cmd/`, `internal/`, `client/`, `assets/`)
+ Implemented a pluggable ABR policy engine using paper-aligned bandwidth formulas (`Tc - Tg`, `Tc - Tn`)
+ Added segment-level in-memory cache tracking (`SegmentCache`)
+ Built throughput estimation module using Exponential Moving Average (EMA) for uncached/cached streams
+ Integrated Prometheus for metrics collection and Jaeger-compatible OpenTelemetry tracing
+ Developed full test suite for cache, estimator, and ABR logic (all unit tests passing ✅)
+ Created a DASH.js-based test client to simulate real adaptive streaming from the browser
+ Set up static routing for MPD + segment files under `/videos`, and served the player UI from `/`
+ Validated successful dynamic MPD rewriting and client-aware segment quality adaptation

---

## Routes

- `-X [PUT] '/api/videos'`: upload a new video
- `-X [GET] '/api/videos'`: get a list of videos with their CID
- `-X [GET] '/api/videos/<cid>'`: get a video MPD by its CID
- `-X [GET] '/api/streams/<cid>'`: handle segment streaming over DASH
- `-X [GET] '/metrics'`: expose Prometheus metrics
- `/`: serves the DASH.js HTML client
- `/videos`: serves static `.mpd` and `.m4s` files for testing
