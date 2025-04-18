#!/bin/bash

# check if MP4Box is installed
if ! command -v MP4Box &> /dev/null; then
    echo "MP4Box could not be found. Please install it before running this script."
    exit 1
fi

# read input video files and output directory from input arguments
if [ "$#" -lt 2 ]; then
    echo "usage: $0 <output_directory> <input_file1> [<input_file2> ...]"
    exit 1
fi

# assign the first argument as the output directory and the rest as input files
output_dir="$1"
shift
input_files="$@"#!/bin/bash

# check if MP4Box is installed
if ! command -v MP4Box &> /dev/null; then
    echo "MP4Box could not be found. Please install it before running this script."
    exit 1
fi

# read input video files and output directory from input arguments
if [ "$#" -lt 2 ]; then
    echo "usage: $0 <output_directory> <input_file1> [<input_file2> ...]"
    exit 1
fi

# assign the first argument as the output directory and the rest as input files
output_dir="$1"
shift
input_files="$@"

# create the output directory if it doesn't exist
mkdir -p "$output_dir"

# run MP4Box with the provided inputs
MP4Box -dash 4000 -frag 4000 -rap -segment-name "$output_dir/segment_" -out "$output_dir/manifest.mpd" $input_files
