#!/bin/bash

# base directory and URLs
BASE_DIR="bp/videos"
URLS=(
    "https://www.sample-videos.com/video321/mp4/720/big_buck_bunny_720p_5mb.mp4"
    "https://www.sample-videos.com/video321/mp4/720/big_buck_bunny_720p_10mb.mp4"
    "https://www.sample-videos.com/video321/mp4/720/big_buck_bunny_720p_20mb.mp4"
    "https://www.sample-videos.com/video321/mp4/720/big_buck_bunny_720p_30mb.mp4"
)

# create the base directory
mkdir -p "$BASE_DIR"

# download each file
for URL in "${URLS[@]}"; do
    FILENAME=$(basename "$URL")
    wget -O "$BASE_DIR/$FILENAME" "$URL"
done

# check if all files are downloaded
ALL_DOWNLOADED=true
for URL in "${URLS[@]}"; do
    FILENAME=$(basename "$URL")
    if [ ! -f "$BASE_DIR/$FILENAME" ]; then
        ALL_DOWNLOADED=false
        break
    fi
done

# output the result
if $ALL_DOWNLOADED; then
    echo "Files downloaded successfully"
else
    echo "Some files were not downloaded"
fi
