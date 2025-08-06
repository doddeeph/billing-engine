# syntax=docker/dockerfile:1
FROM golang:1.24.5-alpine3.22

# Install dependencies
RUN apk add --no-cache wget git bash

# Install wait-for-it.sh
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Install golang-migrate CLI
RUN wget -q https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz && \
    tar -xzf migrate.linux-amd64.tar.gz -C /usr/local/bin && \
    rm migrate.linux-amd64.tar.gz

# Create app directory
WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source
COPY . .

# Build the Go binary
RUN go build -o billing-engine ./cmd

# Copy entrypoint
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Set entrypoint and default command
ENTRYPOINT ["/entrypoint.sh"]
CMD ["./billing-engine"]
