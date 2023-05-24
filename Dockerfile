FROM golang:1.19-alpine as build

# Set shell for this stage with proper pipefall
SHELL ["/bin/sh", "-o", "pipefail", "-c"]

ENV CGO_ENABLED=0

WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the whole app for build purposes
COPY . .

# Build go app
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build \
    -ldflags='-w -s' -a \
    -o bux-wallet ./cmd/main.go

# App final stage
FROM alpine as app

# Set shell for this stage with proper pipefall
SHELL ["/bin/sh", "-o", "pipefail", "-c"]

# Install CA certificates, curl, tzdata and clean up
RUN apk update && \
    apk add --no-cache ca-certificates curl tzdata && \
    rm -rf /var/cache/apk/*

# Copy binary and config files from /build
# to root folder of alpine container.
COPY --from=build ["/build/bux-wallet", "/"]
COPY --from=build ["/build/Procfile", "/"]

# command for debug purposes
CMD ["/bux-wallet"]
