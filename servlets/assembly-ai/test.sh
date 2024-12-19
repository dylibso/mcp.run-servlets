extism call \
  --allow-path "data:/" \
  --allow-host "api.assemblyai.com" \
  --verbose \
  --config ASSEMBLYAI_API_KEY="$ASSEMBLYAI_API_KEY" \
  --wasi \
  --input "{\"params\":{\"name\":\"transcribe\",\"arguments\":{\"audio_path\":\"hello.mp3\"}}}" \
  --log-level debug \
  ./dist/plugin.wasm \
  call