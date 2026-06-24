#!/bin/bash

echo "Starting tests for pokemon-availability..."

# Executa todos os testes do projeto dentro do container
docker run --rm \
  -v "$PWD":/app \
  -w /app \
  golang:1.26.3-alpine \
  sh -c "go test -v ./..."

echo "Tests finished!"
