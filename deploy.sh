#!/usr/bin/env bash

# Variables
IMAGE_NAME="duthweatherstation"
REGISTRY_NAME="duthweatherstationdocker.azurecr.io"
VERSION="v3"

# Build the Docker image
docker build -t ${IMAGE_NAME} .

# Tag the Docker image
docker tag ${IMAGE_NAME} ${REGISTRY_NAME}/${IMAGE_NAME}:${VERSION}

# Push the Docker image to ACR
docker push ${REGISTRY_NAME}/${IMAGE_NAME}:${VERSION}
