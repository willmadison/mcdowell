#!/usr/bin/env bash

set -ev

CONTAINER_NAME=gcr.io/atlblacktech-slack-bot/mcdowell

PROJECT_NAME='github.com/willmadison/mcdowell'
PROJECT_DIR=${PWD}

CONTAINER_PROJECT_ROOT='/root/'
CONTAINER_PROJECT_DIR="${CONTAINER_PROJECT_ROOT}/${PROJECT_NAME}"
BUILD_VERSION=${CIRCLE_BUILD_NUM}.$((CIRCLE_NODE_INDEX + 1))

echo "Current directory contents.... $PROJECT_DIR"

ls -Flah $PROJECT_DIR

docker run --rm \
    --net="host" \
    -v ${PROJECT_DIR}:${CONTAINER_PROJECT_DIR} \
    -e CI=true \
    -w "${CONTAINER_PROJECT_DIR}" \
    golang:1.22.5-alpine \
    go test

docker run --rm \
        --net="host" \
        -v ${PROJECT_DIR}:${CONTAINER_PROJECT_DIR} \
        -e CGO_ENABLED=0 \
        -w "${CONTAINER_PROJECT_DIR}" \
        golang:1.22.5-alpine \
        go build -v -ldflags "-X main.version=${BUILD_VERSION}" ${PROJECT_NAME}/cmd/mcdowell

docker build -f ${PROJECT_DIR}/Dockerfile \
    -t ${CONTAINER_NAME}:${BUILD_VERSION} \
    "${PROJECT_DIR}"

rm -f "${PROJECT_DIR}/mcdowell"

docker tag ${CONTAINER_NAME}:${BUILD_VERSION} ${CONTAINER_NAME}:latest

sudo gcloud docker -- push gcr.io/atlblacktech-slack-bot/mcdowell
sudo kubectl set image deployment/atlblacktech-slack-bot atlblacktech-slack-bot=gcr.io/atlblacktech-slack-bot/mcdowell:${BUILD_VERSION}