build:
	npm run build

test: build
	xtp plugin test ../../servlets/crypto-hash/dist/plugin.wasm --with ./dist/plugin.wasm
