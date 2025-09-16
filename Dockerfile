# ----------------------
# 1️⃣ Build Stage
# ----------------------
FROM golang:1.23 AS builder

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build a static binary for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o clipsync .

# ----------------------
# 2️⃣ Runtime Stage
# ----------------------
FROM alpine:3.20

# Add certificates for HTTPS if your app needs it
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/clipsync .

# Expose the port your app listens on
EXPOSE 9000 

# Run the binary
CMD ["./clipsync"]
