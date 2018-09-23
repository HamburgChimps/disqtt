#!/bin/bash

# Make sure the script runs in the directory in which it is placed
cd $(dirname `[[ $0 = /* ]] && echo "$0" || echo "$PWD/${0#./}"`)

docker build . -t disqtt/alpine-builder:$GOLANG_VERSION --build-arg GOLANG_VERSION=$GOLANG_VERSION
