.PHONY: run clean

# Build target with caching - only rebuilds if source files change
dist/plugin.wasm: $(shell find . -name '*.go')
	xtp plugin build

# Run the WordPress categories list command
run: dist/plugin.wasm
	./test.sh

# Clean built files
clean:
	rm -f dist/plugin.wasm

# Default target
.DEFAULT_GOAL := run