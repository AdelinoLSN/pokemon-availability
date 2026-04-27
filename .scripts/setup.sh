#!/bin/bash

echo "Starting..."

docker run --rm -v "$PWD":/app -w /app golang:1.22-alpine go mod init pokemon-availability || echo "go.mod already exists."

docker run --rm -v "$PWD":/app -w /app golang:1.22-alpine sh -c "go get github.com/lib/pq && go mod tidy"

echo "Files go.mod and go.sum are ready!"
echo "Now you can run: docker-compose up --build"
