FROM golang:1.11-alpine3.7 as shell

ARG APPNAME

# Set working directory
RUN mkdir -p /go/src/$APPNAME
WORKDIR /go/src/$APPNAME

# Install development dependencies
# * curl bash gawk diffutils expect for ts
RUN apk add --no-cache curl bash gawk diffutils expect && \
    curl -o /usr/local/bin/ts -L https://raw.githubusercontent.com/thinkerbot/ts/v2.0.2/bin/ts && \
    chmod +x /usr/local/bin/ts && \
    apk add --no-cache git build-base

# Add build dependencies
RUN go get github.com/spf13/pflag

COPY . .

#############################################################################
FROM shell as build
RUN go build && go install -v

#############################################################################
FROM alpine:3.7 as app
ENV PATH="/app/bin:$PATH"
WORKDIR /app/
RUN apk --no-cache add ca-certificates
ARG APPNAME
COPY --from=build /go/bin/$APPNAME /app/bin/$APPNAME