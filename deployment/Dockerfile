# Build Stage
FROM golang:1.24-alpine AS builder

# Set Go Env
ENV GOPROXY=https://goproxy.io,direct
ENV GO111MODULE=on

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o appsite-monolith ./cmd/appsite-monolith

# Run Stage
FROM alpine:latest

# Certificates and Timezone
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary
COPY --from=builder /app/appsite-monolith .

# Copy configs (User should mount real config in production, but we copy default)
COPY configs/ ./configs/

# Create necessary directories
RUN mkdir -p logs storage/db

EXPOSE 8080

CMD ["./appsite-monolith"]
