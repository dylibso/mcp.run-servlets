#!/bin/bash
set -e

PLUGIN_PATH="target/wasm32-wasip1/release/plugin.wasm"


# Test list_containers
echo "Testing list_containers..."
xtp plugin call "$PLUGIN_PATH" call --wasi --allow-host=localhost --config DOCKER_API_ENDPOINT="http://localhost:2376" -i '{"method": "call", "params": {"name": "list_containers", "arguments": {}}}' > list_containers.out
cat list_containers.out | jq

# Extract the container ID for 'wonderful_lovelace' (if present)
CONTAINER_ID=$(jq -r '.content[0].text' < list_containers.out | jq -r '.[] | select(.Names[] == "/wonderful_lovelace") | .Id // empty')
if [ -z "$CONTAINER_ID" ]; then
  echo "Container 'wonderful_lovelace' not found. Skipping exec_in_container test."
  exit 0
fi

echo '------------'

echo "Testing exec_in_container on container: $CONTAINER_ID ('wonderful_lovelace')..."
xtp plugin call "$PLUGIN_PATH" call --wasi --allow-host=localhost --config DOCKER_API_ENDPOINT="http://localhost:2376" -i '{"method": "call", "params": {"name": "exec_in_container", "arguments": {"container_id": "'$CONTAINER_ID'", "cmd": ["ls", "-la"]}}}' | jq
