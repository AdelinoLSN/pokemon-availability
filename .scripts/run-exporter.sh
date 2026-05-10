#!/bin/sh

set -e

IMAGE_NAME="pokemon-availability-exporter"
NETWORK_NAME="pokemon-availability_default"
ENV_FILE=".env.docker"

echo "Building Exporter Docker image..."

docker build \
  -f Dockerfile.exporter \
  -t ${IMAGE_NAME}:latest \
  .

echo ""
echo "Running Exporter container..."

mkdir -p .outputs

docker run --rm \
  --network ${NETWORK_NAME} \
  --env-file ${ENV_FILE} \
  -v "$(pwd)/.outputs:/app/.outputs" \
  ${IMAGE_NAME}:latest

echo ""
echo "Exporter finished successfully."
