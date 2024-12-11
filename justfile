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
  xtp plugin test --log-level warn

  cd ../qr-code
  xtp plugin test --log-level warn

  cd ../currency-converter
  xtp plugin test --allow-host api.fxratesapi.com --log-level warn

  cd ../eval_js
  xtp plugin test --log-level warn

  cd ../fetch
  xtp plugin test --allow-host getxtp.com --log-level warn