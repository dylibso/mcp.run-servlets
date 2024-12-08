ifeq ($(origin WASI_SDK_PATH),undefined)
    WASI_SDK_PATH=./wasi-sdk
    ifeq ("$(wildcard ./wasi-sdk)","")
        $(shell chmod +x ./install-wasi-sdk.sh)
        $(shell ./install-wasi-sdk.sh)
    endif
endif

.PHONY: build
build: dist/plugin.wasm

.PHONY: fetchdeps
fetchdeps: include/extism-pdk.hpp include/magic_enum/include/magic_enum include/jsoncons/include


.PHONY: format
format:
	$(WASI_SDK_PATH)/bin/clang-format -i *.cpp *.hpp

.PHONY: clean
clean:
	rm -f dist/*.wasm *.o

.PHONY: veryclean
veryclean: clean
	rm -rf include ./wasi-sdk


include/extism-pdk.hpp:
	mkdir -p include
	cd include && \
		curl -OsL https://raw.githubusercontent.com/extism/cpp-pdk/refs/tags/v0.1.0/extism-pdk.hpp

include/magic_enum/include/magic_enum:
	mkdir -p include/magic_enum
	cd include/magic_enum && \
		curl -OsL https://github.com/Neargye/magic_enum/releases/download/v0.9.6/magic_enum-v0.9.6.tar.gz && \
		tar xf magic_enum-v0.9.6.tar.gz

include/jsoncons/include:
	mkdir -p include
	cd include && \
		curl -OsL https://github.com/danielaparker/jsoncons/archive/refs/tags/v0.177.0.zip && \
		unzip -qq v0.177.0.zip && \
		mv jsoncons-0.177.0 jsoncons

pdk.gen.o: pdk.gen.cpp pdk.gen.hpp include/extism-pdk.hpp include/magic_enum/include/magic_enum include/jsoncons/include
	$(WASI_SDK_PATH)/bin/clang++ -Iinclude -Iinclude/jsoncons/include -Iinclude/magic_enum/include/magic_enum -std=c++23 -fno-exceptions -O2 -g -c pdk.gen.cpp

impl.o: impl.cpp pdk.gen.hpp include/extism-pdk.hpp include/magic_enum/include/magic_enum include/jsoncons/include
	$(WASI_SDK_PATH)/bin/clang++ -Iinclude -Iinclude/jsoncons/include -Iinclude/magic_enum/include/magic_enum -std=c++23 -fno-exceptions -O2 -g -c impl.cpp

dist/plugin.wasm: pdk.gen.o impl.o
	mkdir -p dist
	$(WASI_SDK_PATH)/bin/clang++ -std=c++23 -fno-exceptions -O2 -g -o $@ $^ -mexec-model=reactor
