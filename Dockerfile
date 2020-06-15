FROM golang:1.14

# Setup appuser
RUN groupadd -g 1001 appuser && \
    useradd -r -u 1001 -g appuser appuser -m -s /bin/bash

# Install workflow dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates sudo vim less curl jq git man gettext pandoc && \
    adduser appuser sudo && \
    printf "%s\n" "appuser ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

# Set working directory
ARG PACKAGE=pkg
RUN mkdir -p /go/src/$PACKAGE && \
    ln -s /go/src/$PACKAGE/man /go/man
WORKDIR /go/src/$PACKAGE

# Add project dependencies
ENV GO111MODULE=auto
RUN apt-get install -y --no-install-recommends expect && \
    curl -o /usr/local/bin/ts -L https://raw.githubusercontent.com/thinkerbot/ts/v2.0.3/bin/ts && \
    chmod +x /usr/local/bin/ts
USER appuser
COPY --chown=appuser:appuser go.mod go.sum /go/src/$PACKAGE/
RUN go mod download
