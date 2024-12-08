set shell := ["bash", "-euo", "pipefail"]

build:
  #!/usr/bin/env sh
  for dir in servlets/*/; do
    cd "$dir"
    bash ./prepare.sh
    xtp plugin build
    cd ../..
  done

push:
  #!/usr/bin/env sh
  for dir in servlets/*/; do
    cd "$dir"
    bash ./prepare.sh
    xtp plugin push
    cd ../..
  done

