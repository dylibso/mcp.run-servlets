
build:
  #!/usr/bin/env bash
  set -eou pipefail
  for dir in servlets/*/; do
    cd "$dir"
    bash ./prepare.sh
    xtp plugin build
    cd ../..
  done

push:
  #!/usr/bin/env bash
  set -eou pipefail
  for dir in servlets/*/; do
    cd "$dir"
    bash ./prepare.sh
    xtp plugin push
    cd ../..
  done

test:
  #!/usr/bin/env bash
  set -eou pipefail
  xtp plugin build --path test/testsuite

  for dir in servlets/*/; do
    cd "$dir"
    echo "Testing $dir"
    xtp plugin test --allow-host '*' --log-level warn
    cd ../..
  done
