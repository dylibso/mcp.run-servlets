extism call \
  --allow-path "data:/" \
  --allow-host "*" \
  --verbose \
  --config client_id="$GOOGLE_CALENDAR_CLIENT_ID" \
  --config client_secret="$GOOGLE_CALENDAR_CLIENT_SECRET" \
  --wasi \
  --input "{\"params\":{\"name\":\"login-initiate\",\"arguments\":{ \"device_code\": \"\"}}}" \
  --log-level debug \
  ./dist/plugin.wasm \
  call