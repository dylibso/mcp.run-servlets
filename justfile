
[no-cd]
prepare:
    #!/usr/bin/env bash
    set -eou pipefail
    if [ -f "./prepare.sh" ]; then
        bash ./prepare.sh || exit 1
    fi
    if [ -f "./pyproject.toml" ]; then
        uv sync
    fi

build:
  #!/usr/bin/env bash
  set -eou pipefail

  for dir in servlets/*/; do
    cd "$dir"
    echo "Building $dir"
    just prepare
    xtp plugin build
    cd ../..
  done

  cd simulations/describe-output
  make build

push:
  #!/usr/bin/env bash
  set -eou pipefail
  for dir in servlets/*/; do
    cd "$dir"
    just prepare
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
