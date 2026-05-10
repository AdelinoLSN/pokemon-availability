#!/bin/sh

set -e

IMAGE_NAME="pokemon-availability-etl"
NETWORK_NAME="pokemon-availability_default"
ENV_FILE=".env.docker"

echo "Building ETL Docker image..."

docker build \
  -f Dockerfile.etl \
  -t ${IMAGE_NAME}:latest \
  .

echo ""
echo "Running ETL container..."

docker run --rm \
  --network ${NETWORK_NAME} \
  --env-file ${ENV_FILE} \
  ${IMAGE_NAME}:latest

echo ""
echo "ETL finished successfully."
