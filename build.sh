#!/bin/bash

APP_NAME="markdown-server"
OUTPUT_DIR="./bin"
PLATFORMS=("darwin/arm64" "linux/arm64" "darwin/amd64" "linux/amd64" "windows/amd64")

# Create the output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Compile for each platform
for PLATFORM in "${PLATFORMS[@]}"; do
  OS=$(echo $PLATFORM | cut -d'/' -f1)
  ARCH=$(echo $PLATFORM | cut -d'/' -f2)
  OUTPUT_NAME="$OUTPUT_DIR/$APP_NAME-$OS-$ARCH"
  if [ "$OS" = "windows" ]; then
    OUTPUT_NAME+='.exe'
  fi

  echo "Compiling for $OS/$ARCH..."

  # Set the GOOS and GOARCH environment variables and build the binary
  env GOOS=$OS GOARCH=$ARCH go build -o "$OUTPUT_NAME"

  if [ $? -ne 0 ]; then
    echo "An error occurred while compiling for $OS/$ARCH"
    exit 1
  fi
done

echo "Binaries successfully compiled!"
