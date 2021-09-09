#!/usr/bin/env bash

# installing postgres helm chart under a given release name
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add kube-state-metrics https://kubernetes.github.io/kube-state-metrics
helm repo update

# install jaeger dependency
helm upgrade --install telemetry ./charts/telemetry
helm upgrade --install  prometheus prometheus-community/prometheus

# install database helm charts for service
helm upgrade --install merchant-component-db -f ./kubernetes/merchant-component-db/values.yaml bitnami/postgresql
helm upgrade --install shopper-component-db -f ./kubernetes/shopper-component-db/values.yaml bitnami/postgresql
helm upgrade --install auth-service-db -f ./kubernetes/auth-service/postgresql/values.yaml bitnami/postgresql

# install redis helm charts for auth service
helm upgrade --install auth-service-redis -f ./kubernetes/auth-service/redis/values.yaml bitnami/redis

# install authentication service helm chart
helm upgrade --install auth-service ./charts/authentication_service
