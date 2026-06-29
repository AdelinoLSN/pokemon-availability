#!/bin/bash

set -e

echo "Starting coverage analysis for pokemon-availability..."

MIN_COVERAGE="${MIN_COVERAGE:-80.0}"

# 1. Roda os testes gerando o arquivo bruto de profile (coverage.out)
docker run --rm \
  -v "$PWD":/app \
  -w /app \
  golang:1.26.3-alpine \
  sh -c "go test -coverprofile=coverage.out ./..."

# 2. Valida a cobertura mínima total
COVERAGE=$(
  docker run --rm \
    -v "$PWD":/app \
    -w /app \
    golang:1.26.3-alpine \
    sh -c "go tool cover -func=coverage.out" |
    awk '/^total:/ {gsub(/%/, "", $3); print $3}'
)

awk -v coverage="$COVERAGE" -v minimum="$MIN_COVERAGE" 'BEGIN {
  if (coverage + 0 < minimum + 0) {
    printf "Coverage %.1f%% is below required %.1f%%\n", coverage, minimum
    exit 1
  }

  printf "Coverage %.1f%% meets required %.1f%%\n", coverage, minimum
}'

# 3. Converte o profile bruto em um arquivo HTML (coverage.html) legível
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
