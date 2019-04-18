#!/bin/bash

VERSION=0.0.5

DIR="$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")"
cd "${DIR}"

gofmt -d pkg cmd
gofmt -w pkg cmd

operator-sdk generate k8s

operator-sdk build eu.gcr.io/infra-sandbox-58fe57e9/masterservice-operator:${VERSION}
docker push eu.gcr.io/infra-sandbox-58fe57e9/masterservice-operator:${VERSION}

kubectl apply -f deploy/crds/blablacar_v1_masterservice_crd.yaml
kubectl apply -f deploy/service_account.yaml
kubectl apply -f deploy/role.yaml
kubectl apply -f deploy/role_binding.yaml
sed -e 's|REPLACE_IMAGE|eu.gcr.io/infra-sandbox-58fe57e9/masterservice-operator:'"${VERSION}"'|g' deploy/operator.yaml | kubectl apply -f -
