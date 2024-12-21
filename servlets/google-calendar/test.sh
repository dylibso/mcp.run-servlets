extism call \
  --allow-path "data:/" \
  --allow-host "*" \
  --verbose \
  --config client_id="$GOOGLE_CALENDAR_CLIENT_ID" \
  --config client_secret="$GOOGLE_CALENDAR_CLIENT_SECRET" \
  --wasi \
  --input "{\"params\":{\"name\":\"google-calendar-login\",\"arguments\":{ \"device_code\": \"AH-1Ng3dDp2_HBT50fyB1Lh8D5-y36XVyq4dYtpZAzo1BZzVt0GP6W18CuVUynTAJt5RxAIDpI0Hgze-piNqcDaApKaZRNocLg\"}}}" \
  --log-level debug \
  ./dist/plugin.wasm \
  call