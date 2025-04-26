#!/bin/bash

scripts/videos/fetch.sh
scripts/videos/encode.sh

docker compose up -d ipfs0 ipfs1 ipfs2

scripts/cluster/bootstrap.sh

docker compose up -d bootstrap

docker compose up -d telescope jaeger prometheus
