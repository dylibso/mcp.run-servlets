#!/bin/bash

# Get API key from environment
API_KEY=${TENOR_API_KEY:-""}

if [ -z "$API_KEY" ]; then
    echo "Error: TENOR_API_KEY environment variable not set"
    exit 1
fi

# Replace placeholder in source code
sed -i.bak 's/defaultAPIKey := "TENOR_API_KEY_PLACEHOLDER"/defaultAPIKey := "'$API_KEY'"/' main.go

# Build the plugin
tinygo build -target wasi -o dist/plugin.wasm .

# Restore original source
mv main.go.bak main.go