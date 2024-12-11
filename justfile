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

  cd servlets/greet
  xtp plugin test --log-level debug