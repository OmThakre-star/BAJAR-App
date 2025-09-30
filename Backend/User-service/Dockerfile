# Start from the official Golang image for building
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd/server/main.go

# Start a minimal image for running
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/user-service .
COPY config/userConfig.yaml ./config/userConfig.yaml
EXPOSE 8080
CMD ["./user-service"]
