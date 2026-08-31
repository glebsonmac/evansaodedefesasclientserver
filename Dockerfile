# Multi-stage Dockerfile for d3c server
FROM golang:1.21-alpine AS builder
WORKDIR /src

# Cache modules
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org && go mod download

# Copy full project and build the server
COPY . .
WORKDIR /src/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /dist/d3c-server

# Final image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /dist/d3c-server /usr/local/bin/d3c-server
WORKDIR /app

# Default environment (can be overridden at runtime)
ENV D3C_SERVIDOR=0.0.0.0
ENV D3C_PORTA=80
ENV D3C_CONEXAO=https
ENV LISTENER_TYPE=http
ENV LISTENER_PORT=80

EXPOSE 80
ENTRYPOINT ["/usr/local/bin/d3c-server"]
