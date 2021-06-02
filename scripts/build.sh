#!/usr/bin/env bash

set -euo pipefail

GOOS="linux" go build -ldflags='-s -w' -o bin/helper github.com/paketo-buildpacks/libjvm/cmd/helper
GOOS="linux" go build -ldflags='-s -w' -o bin/main github.com/paketo-buildpacks/bellsoft-liberica/cmd/main

if [ "${STRIP:-false}" != "false" ]; then
  strip bin/helper bin/main
fi

if [ "${COMPRESS:-false}" != "false" ]; then
  upx -q -9 bin/helper bin/main
fi

ln -fs main bin/build
ln -fs main bin/detect
