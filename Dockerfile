# Specify the Go version
ARG GO_VERSION=1
FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /usr/src/app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the entire source and build the Go binary
COPY . .
RUN go build -v -o /run-app ./cmd/api

# Use a smaller base image
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /run-app /usr/local/bin/

# Ensure proper handling of dependencies
RUN apk --no-cache add ca-certificates

# Set the binary as the entrypoint
CMD ["run-app"]
