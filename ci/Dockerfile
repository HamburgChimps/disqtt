ARG GOLANG_VERSION=1.11
FROM golang:${GOLANG_VERSION}


RUN apk add --update --no-cache \
git bash openssh openssl python build-base cmake libc6-compat linux-headers \
    && rm -rf /tmp/* \
    && adduser -D -u 1000 disqtt