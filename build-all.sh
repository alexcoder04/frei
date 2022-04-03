#!/bin/sh

for arch in amd64 386 arm; do
  echo "build $arch..."
  export GOOS="linux"
  export GOARCH="$arch"
  export VERSION="$(git describe --tags --abbrev=0)"
  export COMMIT="$(git rev-list -1 HEAD)"
  go build \
    -ldflags "-X main.Version=$VERSION -X main.CommitSHA=$COMMIT" \
    -o "build/frei-linux-$arch" .
done

