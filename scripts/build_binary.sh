#!/usr/bin/env bash

set -euo pipefail

OS="$(uname -s)"
ARCH="$(uname -m)"

case $OS in
    "Linux")
        case $ARCH in
        "x86_64")
            ARCH=amd64
            ;;
        "armv6")
            ARCH=armv6l
            ;;
        "armv8")
            ARCH=arm64
            ;;
        .*386.*)
            ARCH=386
            ;;
        esac
        PLATFORM="linux"
    ;;
    "Darwin")
        PLATFORM="darwin"
        ARCH=amd64
    ;;
esac

echo 'Building binary'
env GOOS=$PLATFORM GOARCH=$ARCH go build -tags 'bindatafs'

if [[ $? -ne 0 ]]; then
    echo 'An error has occurred! Aborting the script execution...'
    exit 1
fi
echo 'Finished building binary'
