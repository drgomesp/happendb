FROM golang:1.18-alpine AS build-env

# Install minimum necessary dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3
RUN apk add --no-cache $PACKAGES

# Set working directory for the build
WORKDIR /happendb

# Add source files
COPY . .

RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN task build

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates

# Copy over binaries from the build-env
COPY --from=build-env /happendb/happendbd /usr/bin/happendbd

EXPOSE 26656 26657 1317 9090

# Run happendbd by default
CMD ["happendbd"]