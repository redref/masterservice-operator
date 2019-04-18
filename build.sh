#!/bin/bash

DIR="$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")"
cd "${DIR}"

gofmt -d pkg cmd
gofmt -w pkg cmd

operator-sdk generate k8s
