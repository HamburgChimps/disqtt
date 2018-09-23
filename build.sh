#!/bin/bash
GOLANG_VERSION=1.11-rc-alpine
BASE_IMAGE="disqtt/alpine-builder:$GOLANG_VERSION"

# Make sure the script runs in the directory in which it is placed
cd $(dirname `[[ $0 = /* ]] && echo "$0" || echo "$PWD/${0#./}"`)

GOLANG_VERSION=$GOLANG_VERSION ci/ensure.sh

docker run --rm \
-v "$(pwd):/disqtt" \
-v "$(pwd)/vendor:/go/pkg" \
-w /disqtt \
$BASE_IMAGE ash -c \
"echo 'Building binaries...' \
&& CGO_ENABLED=0 go build -o ./bin/client ./cmd/client  \
&& CGO_ENABLED=0 go build -o ./bin/service ./cmd/service  \
&& chown -R $(id -u):$(id -g) /go/pkg ./bin/*"