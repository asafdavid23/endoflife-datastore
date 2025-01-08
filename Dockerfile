FROM golang:1.23-alpine as builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o datastore ./cmd/main.go

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Install CA certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /root/

# Copy the compiled binary and configuration file from the builder stage
COPY --from=builder /app/datastore .
COPY config/config.yaml ./config.yaml

# Expose the port for the HTTP API
EXPOSE 8080

# Set the entrypoint command
ENTRYPOINT ["./datastore"]