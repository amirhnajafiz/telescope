#!/bin/bash

# input and output directories
INPUT_DIR="bp/videos"
OUTPUT_DIR="bp/idp"

# create the output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# process each video in the input directory
for VIDEO in "$INPUT_DIR"/*.mp4; do
    # extract the base name of the video (without extension)
    BASENAME=$(basename "$VIDEO" .mp4)
    
    # create a directory for the output of this video
    VIDEO_OUTPUT_DIR="$OUTPUT_DIR/$BASENAME"
    mkdir -p "$VIDEO_OUTPUT_DIR"
    
    # run the ffmpeg command
    ffmpeg -i "$VIDEO" \
        -map 0:v -b:v:0 300k -s:v:0 426x240 \
        -map 0:v -b:v:1 700k -s:v:1 640x360 \
        -map 0:v -b:v:2 1500k -s:v:2 854x480 \
        -map 0:v -b:v:3 3000k -s:v:3 1280x720 \
        -map 0:a -b:a:0 128k \
        -c:v libx264 -c:a aac -ar 48000 -ac 2 \
        -f dash -seg_duration 4 \
        -use_timeline 1 -use_template 1 \
        -adaptation_sets "id=0,streams=v id=1,streams=a" \
        "$VIDEO_OUTPUT_DIR/stream.mpd"
done

echo "Encoding completed for all videos in $INPUT_DIR."
