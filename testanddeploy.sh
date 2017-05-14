#!/usr/bin/env bash

set -ev

CONTAINER_NAME=gcr.io/atlblacktech-slack-bot/mcdowell

PROJECT_NAME='github.com/willmadison/mcdowell'
PROJECT_DIR=${PWD}

CONTAINER_GOPATH='/go'
CONTAINER_PROJECT_DIR="${CONTAINER_GOPATH}/src/${PROJECT_NAME}"
CONTAINER_PROJECT_GOPATH="${CONTAINER_GOPATH}"

docker run --rm \
    --net="host" \
    -v ${PROJECT_DIR}:${CONTAINER_PROJECT_DIR} \
    -e CI=true \
    -e GOPATH=${CONTAINER_PROJECT_GOPATH} \
    -w "${CONTAINER_PROJECT_DIR}" \
    golang:1.8.1-alpine \
    go test

docker run --rm \
        --net="host" \
        -v ${PROJECT_DIR}:${CONTAINER_PROJECT_DIR} \
        -e GOPATH=${CONTAINER_PROJECT_GOPATH} \
        -e CGO_ENABLED=0 \
        -w "${CONTAINER_PROJECT_DIR}" \
        golang:1.8.1-alpine \
        go build -v ${PROJECT_NAME}/cmd/mcdowell

docker build -f ${PROJECT_DIR}/Dockerfile \
    -t ${CONTAINER_NAME}:latest \
    "${PROJECT_DIR}"

rm -f "${PROJECT_DIR}/mcdowell"

gcloud docker -- push gcr.io/atlblacktech-slack-bot/mcdowell
kubectl set image deployment/atlblacktech-slack-bot atlblacktech-slack-bot=gcr.io/atlblacktech-slack-bot/mcdowell:latest