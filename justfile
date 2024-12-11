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

test:
  #!/usr/bin/env sh
  xtp plugin build --path test/host
  xtp plugin build --path test/testsuite

  for dir in servlets/*/; do
    cd "$dir"
    echo "Testing $dir"
    xtp plugin test --allow-host api.fxratesapi.com --allow-host getxtp.com --log-level warn
    cd ../..
  done