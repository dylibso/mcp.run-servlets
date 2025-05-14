#!/bin/bash
set -e

PLUGIN_PATH="target/wasm32-wasip1/release/plugin.wasm"

# Test create_minions
echo "Testing create_minions..."
xtp plugin call "$PLUGIN_PATH" call --wasi -i '{"method": "call", "params": {"name": "create_minions", "arguments": {"prompts": ["do A", "do B", "do C"]}}}' \
    --allow-host 'www.mcp.run' --allow-host 'api.cloudflare.com' \
    --config 'MINION_ENDPOINT=https://www.mcp.run/api/runs/evacchi/demo-profile/minions?nonce=9ZwWd-OGN7OMmMlVNi87Ng&sig=KymBobHhfR9gkwBFsXlYVYEBCoNpS3PI7F8pAnZ5DVc' \
    --config CF_ACCOUNT_ID=92e6c721ef19db523a0d9de0265906df --config CF_DATABASE_ID=bdf722a8-46ab-4e6a-b86b-d26be1cc4b96 --config CF_API_TOKEN=obbqy4I57SxxQVVjUYzMmKeQqttUNhDcYQ6NQ4At > create_minions.out
cat create_minions.out

# Extract a minion_id and mob_id for further tests
MINION_ID=$(jq -r '.content[0].text' < create_minions.out | jq -r '.minion_ids[0]')
MOB_ID=$(jq -r '.content[0].text' < create_minions.out | jq -r '.mob_id')

echo "\nTesting check_minion_state for minion_id: $MINION_ID..."
echo xtp plugin call "$PLUGIN_PATH" call --wasi -i '{"method": "call", "params": {"name": "check_minion_state", "arguments": {"minion_id": "'$MINION_ID'"}}}' \
    --allow-host 'www.mcp.run' --allow-host 'api.cloudflare.com' \
    --config 'MINION_ENDPOINT=https://www.mcp.run/api/runs/evacchi/demo-profile/minions?nonce=9ZwWd-OGN7OMmMlVNi87Ng&sig=KymBobHhfR9gkwBFsXlYVYEBCoNpS3PI7F8pAnZ5DVc' \
    --config CF_ACCOUNT_ID=92e6c721ef19db523a0d9de0265906df --config CF_DATABASE_ID=bdf722a8-46ab-4e6a-b86b-d26be1cc4b96 --config CF_API_TOKEN=obbqy4I57SxxQVVjUYzMmKeQqttUNhDcYQ6NQ4At > create_minions.out


xtp plugin call "$PLUGIN_PATH" call --wasi -i '{"method": "call", "params": {"name": "check_minion_state", "arguments": {"minion_id": "'$MINION_ID'"}}}' \
    --allow-host 'www.mcp.run' --allow-host 'api.cloudflare.com' \
    --config 'MINION_ENDPOINT=https://www.mcp.run/api/runs/evacchi/demo-profile/minions?nonce=9ZwWd-OGN7OMmMlVNi87Ng&sig=KymBobHhfR9gkwBFsXlYVYEBCoNpS3PI7F8pAnZ5DVc' \
    --config CF_ACCOUNT_ID=92e6c721ef19db523a0d9de0265906df --config CF_DATABASE_ID=bdf722a8-46ab-4e6a-b86b-d26be1cc4b96 --config CF_API_TOKEN=obbqy4I57SxxQVVjUYzMmKeQqttUNhDcYQ6NQ4At > create_minions.out

cat create_minions.out
