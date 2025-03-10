#!/usr/bin/env bash
# Orignal from: https://github.com/Shopify/javy

set -euo pipefail

PATH_TO_SDK="./wasi-sdk"
if [[ ! -d $PATH_TO_SDK ]]; then
    TMPGZ=$(mktemp)
    VERSION_MAJOR="24"
    VERSION_MINOR="0"
    ARCH=$(uname -m)
    if [[ "$(uname -s)" == "Darwin" ]]; then
        curl --fail --location --silent https://github.com/WebAssembly/wasi-sdk/releases/download/wasi-sdk-${VERSION_MAJOR}/wasi-sdk-${VERSION_MAJOR}.${VERSION_MINOR}-${ARCH}-macos.tar.gz --output $TMPGZ
    else
        curl --fail --location --silent https://github.com/WebAssembly/wasi-sdk/releases/download/wasi-sdk-${VERSION_MAJOR}/wasi-sdk-${VERSION_MAJOR}.${VERSION_MINOR}-${ARCH}-linux.tar.gz --output $TMPGZ
    fi
    mkdir $PATH_TO_SDK
    tar xf $TMPGZ -C $PATH_TO_SDK --strip-components=1
fi
