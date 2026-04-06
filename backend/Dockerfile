# ============================================
# Stage 1: Build the Go binary
# ============================================
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# ============================================
# Stage 2: Minimal production image
# ============================================
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl

# Set timezone
ENV TZ=Asia/Jakarta

WORKDIR /app

COPY --from=builder /app/server .

# Create directories for logs and uploads
RUN mkdir -p /app/logs /app/uploads

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["./server"]
