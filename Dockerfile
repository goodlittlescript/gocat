FROM golang:1.13 as base

# Setup app dir and user
RUN apt-get update && \
    mkdir -p /app && \
    groupadd -g 900 appuser && \
    useradd -r -u 900 -g appuser appuser -m -s /bin/bash && \
    chown -R appuser:appuser /app
ENV PATH="/app/bin:$PATH"
WORKDIR /app
USER appuser

#############################################################################
FROM base as shell

USER root
RUN apt-get install -y --no-install-recommends ca-certificates sudo vim less build-essential curl git man expect && \
    adduser appuser sudo && \
    printf "%s\n" "appuser ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers && \
    curl -o /usr/local/bin/ts -L https://raw.githubusercontent.com/thinkerbot/ts/v2.0.2/bin/ts && \
    chmod +x /usr/local/bin/ts
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

COPY . .

#############################################################################
FROM shell as build
RUN go install

#############################################################################
FROM base as app
COPY --from=build /go/bin /app/bin
