#!/usr/bin/env bash

set -euo pipefail

if [[ -d ../go-cache ]]; then
  GOPATH=$(realpath ../go-cache)
  export GOPATH
fi

GOOS="linux" go build -ldflags='-s -w' -o bin/class-counter github.com/paketo-buildpacks/libjvm/cmd/class-counter
GOOS="linux" go build -ldflags='-s -w' -o bin/link-local-dns github.com/paketo-buildpacks/libjvm/cmd/link-local-dns
GOOS="linux" go build -ldflags='-s -w' -o bin/main github.com/paketo-buildpacks/bellsoft-liberica/cmd/main
GOOS="linux" go build -ldflags='-s -w' -o bin/openssl-certificate-loader github.com/paketo-buildpacks/libjvm/cmd/openssl-certificate-loader
GOOS="linux" go build -ldflags='-s -w' -o bin/security-providers-configurer github.com/paketo-buildpacks/libjvm/cmd/security-providers-configurer
ln -fs main bin/build
ln -fs main bin/detect
