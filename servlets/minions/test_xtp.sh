#!/bin/bash
set -e

PLUGIN_PATH="target/wasm32-wasip1/release/plugin.wasm"

# Test create_minions
echo "Testing create_minions..."
xtp plugin call "$PLUGIN_PATH" call --wasi -i '{"method": "call", "params": {"name": "create_minions", "arguments": {"prompts": ["do A", "do B", "do C"]}}}' > create_minions.out
cat create_minions.out

# Extract a minion_id and mob_id for further tests
MINION_ID=$(jq -r '.content[0].text' < create_minions.out | jq -r '.minion_ids[0]')
MOB_ID=$(jq -r '.content[0].text' < create_minions.out | jq -r '.mob_id')

echo "\nTesting check_minion_state for minion_id: $MINION_ID..."
xtp plugin call "$PLUGIN_PATH" call --wasi -i '{"method": "call", "params": {"name": "check_minion_state", "arguments": {"minion_id": "'$MINION_ID'"}}}'

echo "\nTesting check_mob_state for mob_id: $MOB_ID..."
xtp plugin call "$PLUGIN_PATH" call --wasi -i '{"method": "call", "params": {"name": "check_mob_state", "arguments": {"mob_id": "'$MOB_ID'"}}}'
