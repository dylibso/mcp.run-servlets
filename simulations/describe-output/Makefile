build:
	tinygo build -target wasi -o plugin.wasm main.go

test:
	extism call plugin.wasm greet --input "world" --wasi
