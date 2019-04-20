#!/bin/bash

VERSION=0.0.13

DIR="$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")"
cd "${DIR}"

gofmt -d pkg cmd
gofmt -w pkg cmd

operator-sdk generate k8s

operator-sdk build eu.gcr.io/infra-sandbox-58fe57e9/masterservice-operator:${VERSION}
docker push eu.gcr.io/infra-sandbox-58fe57e9/masterservice-operator:${VERSION}

helm upgrade masterservice-operator ./helm --set version=${VERSION}
echo "Waiting start"
sleep 20
helm test masterservice-operator --cleanup
