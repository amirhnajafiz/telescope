#!/bin/bash

# This script connects all IPFS nodes in the cluster to each other.

# function to get the node ID and multiaddresses of a given IPFS container
get_node_info() {
    local container_name=$1
    local node_info=$(docker exec telescope-"$container_name"-1 ipfs id)
    local node_id=$(echo "$node_info" | jq -r '.ID')
    local node_addresses=$(echo "$node_info" | jq -r '.Addresses[]' | grep -v '127.0.0.1' | grep '^/ip4' | head -n 1)
    echo "$node_id" "$node_addresses"
}

# function to connect one IPFS node to another
connect_nodes() {
    local source_container=$1
    local target_address=$2
    local target_id=$3
    docker exec telescope-"$source_container"-1 ipfs swarm connect "$target_address/p2p/$target_id"
}

# get node info for ipfs0, ipfs1, and ipfs2
read ipfs0_id ipfs0_address <<< $(get_node_info ipfs0)
read ipfs1_id ipfs1_address <<< $(get_node_info ipfs1)
read ipfs2_id ipfs2_address <<< $(get_node_info ipfs2)

# connect ipfs0 to ipfs1 and ipfs2
connect_nodes ipfs0 "$ipfs1_address" "$ipfs1_id"
connect_nodes ipfs0 "$ipfs2_address" "$ipfs2_id"

# connect ipfs1 to ipfs0 and ipfs2
connect_nodes ipfs1 "$ipfs0_address" "$ipfs0_id"
connect_nodes ipfs1 "$ipfs2_address" "$ipfs2_id"

# connect ipfs2 to ipfs0 and ipfs1
connect_nodes ipfs2 "$ipfs0_address" "$ipfs0_id"
connect_nodes ipfs2 "$ipfs1_address" "$ipfs1_id"

echo "All IPFS nodes have been connected successfully."
