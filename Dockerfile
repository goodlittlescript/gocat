FROM golang:1.13

# Setup app dir and user
RUN apt-get update && \
    mkdir -p /app && \
    groupadd -g 900 appuser && \
    useradd -r -u 900 -g appuser appuser -m -s /bin/bash && \
    chown -R appuser:appuser /app
ENV PATH="/app/bin:$PATH"
WORKDIR /app

RUN apt-get install -y --no-install-recommends ca-certificates sudo vim less curl jq git man expect && \
    adduser appuser sudo && \
    printf "%s\n" "appuser ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers && \
    curl -o /usr/local/bin/ts -L https://raw.githubusercontent.com/thinkerbot/ts/v2.0.3/bin/ts && \
    chmod +x /usr/local/bin/ts && \
    curl -s -L -O https://github.com/goreleaser/goreleaser/releases/download/v0.135.0/goreleaser_Linux_x86_64.tar.gz && \
    tar -xf goreleaser_Linux_x86_64.tar.gz && rm goreleaser_Linux_x86_64.tar.gz && \
    mv goreleaser /usr/local/bin/goreleaser
USER appuser

# Enable go modules
ENV GO111MODULE=auto

# Set working directory
ARG PACKAGE
RUN mkdir -p /go/src/$PACKAGE
WORKDIR /go/src/$PACKAGE

# Add project dependencies
COPY go.mod go.sum /go/src/$PACKAGE/
RUN go mod download
