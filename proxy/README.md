# Telescope Proxy

The purpose of this proxy is to serve as a gateway for uploading and streaming videos over IPFS. It is a web server implemented using the Golang Fiber framework. Since the proxy server is stateless, it can be easily scaled to accommodate varying workloads. Caching is handled through HTTP headers on the client side, ensuring efficient data retrieval.  

This gateway proxy enhances the efficiency and performance of Adaptive Bitrate (ABR) streaming protocols. Additionally, as a gateway, it provides performance metrics and tracing capabilities, facilitating better monitoring and observability.  

The original Telescope proxy is available at [Telescope Repository](https://github.com/SBUNetSys/Telescope). This project is a rebuilt version of Telescope, designed to address its existing issues while introducing new features and improvements.

## Routes

- `-X [PUT] '/api/videos'`: uploading a new video
- `-X [GET] '/api/videos'`: get a list of videos with their CID
- `-X [GET] '/api/videos/<cid>'`: get a video MPD by its CID
- `-X [GET] '/api/streams/<cid>'`: handling streams over DASH
- `-X [GET] '/metrics'`: expose prometheus metrics
