#!/bin/bash

PROJECT_NAME="github.com/AdelinoLSN/pokemon-availability"

echo "Starting..."

docker run --rm -v "$PWD":/app -w /app golang:1.26-alpine go mod init "$PROJECT_NAME" || echo "go.mod already exists."

docker run --rm -v "$PWD":/app -w /app golang:1.26-alpine sh -c "go mod tidy"

echo "Files go.mod and go.sum are ready!"
echo "Now you can run: docker-compose up --build"
