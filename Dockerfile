# Build go
FROM golang:1.24.4-alpine AS builder
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Set build environment
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Build the application
RUN go build -v -o XrayR .

# Release
FROM alpine:latest
# 安装必要的工具包
RUN apk --update --no-cache add tzdata ca-certificates \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir /etc/XrayR/
COPY --from=builder /app/XrayR /usr/local/bin

ENTRYPOINT [ "XrayR", "--config", "/etc/XrayR/config/config.yml"]