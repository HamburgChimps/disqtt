#!/bin/bash

cd $(dirname `[[ $0 = /* ]] && echo "$0" || echo "$PWD/${0#./}"`)
GOLANG_VERSION=${GOLANG_VERSION:='1.11'}

if docker inspect disqtt/alpine-builder:$GOLANG_VERSION > /dev/null 2>&1 ; then
  echo "GOLANG builder image with version $GOLANG_VERSION exists"
else
  echo "GOLANG builder image with version $GOLANG_VERSION not existing, building it"
  GOLANG_VERSION=$GOLANG_VERSION ./build.sh
fi
