FROM golang:1.11-alpine3.7 as shell

ARG APPNAME

# Set working directory
RUN mkdir -p /go/src/$APPNAME
WORKDIR /go/src/$APPNAME

# Install development dependencies
# * ts: curl bash gawk diffutils expect
# * go: git build-base
RUN apk add --no-cache curl bash gawk diffutils expect && \
    curl -o /usr/local/bin/ts -L https://raw.githubusercontent.com/thinkerbot/ts/v2.0.2/bin/ts && \
    chmod +x /usr/local/bin/ts && \
    apk add --no-cache git build-base

# Enable go modules
ENV GO111MODULE=on

# Add project dependencies
COPY go.mod go.sum /go/src/$APPNAME/
RUN go mod download

COPY . .

#############################################################################
FROM shell as build
RUN go build ./... && go install -v

#############################################################################
FROM alpine:3.7 as app
ENV PATH="/app/bin:$PATH"
WORKDIR /app/
RUN apk --no-cache add ca-certificates
ARG APPNAME
COPY --from=build /go/bin/$APPNAME /app/bin/$APPNAME
