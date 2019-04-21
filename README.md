# Masterservice-operator

[![Go Report Card](https://goreportcard.com/badge/github.com/Junonogis/masterservice-operator)](https://goreportcard.com/report/github.com/Junonogis/masterservice-operator)
[![Build Status](https://travis-ci.org/Junonogis/masterservice-operator.svg?branch=master)](https://travis-ci.org/Junonogis/masterservice-operator)

Masterservice-operator is a tooling operator using a CustomResource named `MasterService`. This custom resource creates 2 services :
  * `name`-all : kubernetes service created with the given `serviceSpec`
  * `name` : kubernetes empty service populated by the operator

The second service will be populated with the oldest `Pod` (based on `StartTime`) found in the first one.

A `callback` field in the `MasterService` definition allow to notify (via HTTP) the chosen `Pod` of his election.

## Usage

This operator is useful to implement a master/slave database cluster with no switchover downtime.

Examples :
  * Mysql Galera Cluster
  * PostgreSQL synchronous cluster
  * Redis

## MasterService example

```
---
# Source: redis/templates/services.yaml
apiVersion: v1
kind: List
items:
  - apiVersion: blablacar.com/v1
    kind: MasterService
    metadata:
      name: release-name
      labels:
        app: redis
        chart: redis-0.1
        release: release-name
        heritage: Tiller
        component: database
    spec:
      serviceSpec:
        type: ClusterIP
        ports:
          - port: 6379
            targetPort: redis
            protocol: TCP
            name: redis
        selector:
          app: redis
          release: release-name
      callback:
        port: 8000
```