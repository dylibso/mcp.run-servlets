#!/bin/bash

# Function to check if a command exists
command_exists () {
  command -v "$1" >/dev/null 2>&1
}

missing_deps=0

# Check for Go
if ! (command_exists go); then
  missing_deps=1
  echo "❌ Go is not installed."
  echo ""
  echo "To install Go, visit the official download page:"
  echo "👉 https://go.dev/dl/"
  echo ""
  echo "Or install it using a package manager:"
  echo ""
  echo "🔹 macOS (Homebrew):"
  echo "    brew install go"
  echo ""
  echo "🔹 Ubuntu/Debian:"
  echo "    sudo apt-get -y install golang-go"
  echo ""
  echo "🔹 Arch Linux:"
  echo "    sudo pacman -S go"
  echo ""
fi

ARCH=$(arch)

# Check for TinyGo
if ! (command_exists tinygo); then
  missing_deps=1
  echo "❌ TinyGo is not installed."
  echo ""
  echo "To install TinyGo, visit the official download page:"
  echo "👉 https://tinygo.org/getting-started/install/"
  echo ""
  echo "Or install it using a package manager:"
  echo ""
  echo "🔹 macOS (Homebrew):"
  echo "    brew tap tinygo-org/tools"
  echo "    brew install tinygo"
  echo ""
  echo "🔹 Ubuntu/Debian:"
  echo "    wget https://github.com/tinygo-org/tinygo/releases/download/v0.31.2/tinygo_0.31.2_$ARCH.deb"
  echo "    sudo dpkg -i tinygo_0.31.2_$ARCH.deb"
  echo ""
  echo "🔹 Arch Linux:"
  echo "    pacman -S extra/tinygo"
  echo ""
fi