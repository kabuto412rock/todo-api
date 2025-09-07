# Go development environment
FROM golang:1.25-bookworm

# Install necessary tools
RUN apt-get update && apt-get install -y \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Set up Go environment
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

# Create workspace directory
WORKDIR /workspace
