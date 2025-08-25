# ---- Build stage ----
FROM golang:1.25.0 AS builder

WORKDIR /app

# Install dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source
COPY backend/ .

# Build binary (disable CGO for small static binary)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o goconda .

# ---- Runtime stage ----
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary only
COPY --from=builder /app/goconda /app/goconda
COPY backend/conf/app.prod.conf /app/conf/app.prod.conf

# Expose port
EXPOSE 8080

# Run app
CMD ["/app/goconda"]
