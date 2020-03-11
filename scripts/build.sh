#!/usr/bin/env bash

set -euo pipefail

if [[ -d ../go-cache ]]; then
  GOPATH=$(realpath ../go-cache)
  export GOPATH
fi

GOOS="linux" go build -ldflags='-s -w' -o bin/build github.com/paketo-buildpacks/bellsoft-liberica/cmd/build
GOOS="linux" go build -ldflags='-s -w' -o bin/class-counter github.com/paketo-buildpacks/libjvm/cmd/class-counter
GOOS="linux" go build -ldflags='-s -w' -o bin/detect github.com/paketo-buildpacks/bellsoft-liberica/cmd/detect
GOOS="linux" go build -ldflags='-s -w' -o bin/link-local-dns github.com/paketo-buildpacks/libjvm/cmd/link-local-dns
GOOS="linux" go build -ldflags='-s -w' -o bin/security-providers-configurer github.com/paketo-buildpacks/libjvm/cmd/security-providers-configurer
