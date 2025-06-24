# --- Build stage ---
FROM golang:1.24-alpine AS builder

ENV PROJECT_PATH=/chirpstack-simulator
ENV CGO_ENABLED=0
ENV GO_EXTRA_BUILD_ARGS="-a -installsuffix cgo"

RUN apk add --no-cache tzdata make git bash

RUN mkdir -p $PROJECT_PATH
WORKDIR $PROJECT_PATH

COPY . .
RUN cd src && go build $GO_EXTRA_BUILD_ARGS -o build/chirpstack-simulator ./main.go

# --- Runtime stage ---
FROM alpine:3.22

ENV PROJECT_PATH=/chirpstack-simulator

WORKDIR $PROJECT_PATH

COPY --from=builder $PROJECT_PATH/src/build/chirpstack-simulator ./chirpstack-simulator

# Set entrypoint (adjust as needed)
ENTRYPOINT ["./chirpstack-simulator"]