# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o file-processor

# Runtime stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/file-processor .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./file-processor"]