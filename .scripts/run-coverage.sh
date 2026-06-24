#!/bin/bash

echo "Starting coverage analysis for pokemon-availability..."

# 1. Roda os testes gerando o arquivo bruto de profile (coverage.out)
docker run --rm \
  -v "$PWD":/app \
  -w /app \
  golang:1.26.3-alpine \
  sh -c "go test -coverprofile=coverage.out ./..."

# 2. Converte o profile bruto em um arquivo HTML (coverage.html) legível
docker run --rm \
  -v "$PWD":/app \
  -w /app \
  golang:1.26.3-alpine \
  sh -c "go tool cover -html=coverage.out -o coverage.html"

echo "Coverage report generated successfully!"
echo "---------------------------------------------------"
echo "To view the report, open the 'coverage.html' file in your browser."
echo "If you are on macOS, you can run: open coverage.html"
echo "---------------------------------------------------"
