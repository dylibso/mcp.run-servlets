#!/bin/bash

# Function to check if a command exists
command_exists () {
  command -v "$1" >/dev/null 2>&1
}

missing_deps=0

# Check for uv
if ! (command_exists uv); then
  missing_deps=1
  echo "❌ uv is not installed."
  echo ""
  echo "To install uv, visit the official download page:"
  echo "👉 https://docs.astral.sh/uv/getting-started/installation/"
  echo ""
fi

# Exit with a bad exit code if any dependencies are missing
if [ "$missing_deps" -ne 0 ]; then
  echo "Install the missing dependencies and ensure they are on your path. Then run this command again."
  # TODO: remove sleep when cli bug is fixed
  sleep 2
  exit 1
fi

# Check for extism-js
if ! command_exists extism-py; then
  echo "❌ extism-py is not installed."
  echo ""
  echo "extism-py is needed to compile the plug-in. You can find the instructions to install it here: https://github.com/extism/python-pdk"
  echo ""
  echo "Alternatively, you can use an install script."
  echo ""
  echo "🔹 Mac / Linux:"
  echo "curl -L https://raw.githubusercontent.com/extism/python-pdk/main/install.sh | bash"
  echo ""
  # TODO: remove sleep when cli bug is fixed
  sleep 2
  exit 1
fi
