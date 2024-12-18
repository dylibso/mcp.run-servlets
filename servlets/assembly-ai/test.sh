xtp plugin build && extism call \
  --allow-host "api.assemblyai.com" \
  --config ASSEMBLYAI_API_KEY="$ASSEMBLYAI_API_KEY" \
  --wasi \
  --input "{\"params\":{\"name\":\"transcribe\",\"arguments\":{\"audio\":\"$(base64 hello.mp3 | tr -d '\n')\"}}}" \
  --log-level debug \
  ./dist/plugin.wasm \
  call