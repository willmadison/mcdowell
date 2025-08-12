#!/usr/bin/env bash
set -euo pipefail

PROJECT_ID="atlblacktech-slack-bot"
REGION="us-central1"
REPO="us.gcr.io"
IMAGE="mcdowell"
BUILD_VERSION="${CIRCLE_BUILD_NUM}.$((CIRCLE_NODE_INDEX + 1))"

PROJECT_NAME='github.com/willmadison/mcdowell'
PROJECT_DIR="${PWD}"
CONTAINER_PROJECT_ROOT='/root'
CONTAINER_PROJECT_DIR="${CONTAINER_PROJECT_ROOT}/${PROJECT_NAME}"

docker run --rm \
  --net="host" \
  -v "${PROJECT_DIR}:${CONTAINER_PROJECT_DIR}" \
  -e CI=true \
  -w "${CONTAINER_PROJECT_DIR}" \
  golang:1.22.5-alpine \
  go test -v

docker run --rm \
  --net="host" \
  -v "${PROJECT_DIR}:${CONTAINER_PROJECT_DIR}" \
  -e CGO_ENABLED=0 \
  -w "${CONTAINER_PROJECT_DIR}" \
  golang:1.22.5-alpine \
  go build -v -ldflags "-X main.version=${BUILD_VERSION}" ${PROJECT_NAME}/cmd/${IMAGE}

gcloud auth configure-docker "${REGION}-docker.pkg.dev"

AR_IMAGE="${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO}/${IMAGE}:${BUILD_VERSION}"
AR_LATEST="${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO}/${IMAGE}:latest"

docker build -f "${PROJECT_DIR}/Dockerfile" -t "${AR_IMAGE}" "${PROJECT_DIR}"
docker tag "${AR_IMAGE}" "${AR_LATEST}"

docker push "${AR_IMAGE}"
docker push "${AR_LATEST}"

kubectl set image deployment/atlblacktech-slack-bot \
  atlblacktech-slack-bot="${AR_IMAGE}"
