#! /usr/bin/env sh

# add jetstack repository
helm repo add jetstack https://charts.jetstack.io || true
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add kube-state-metrics https://kubernetes.github.io/kube-state-metrics
helm repo update

# install cert-manager
helm upgrade --install cert-manager jetstack/cert-manager \
    --set installCRDs=true \
    --namespace default

# wait for cert manager
kubectl rollout status deployment/cert-manager --timeout=2m
kubectl rollout status deployment/cert-manager-webhook --timeout=2m
kubectl rollout status deployment/cert-manager-cainjector --timeout=2m

# install self-signed certificate
cat << 'EOF' | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: self-signed
spec:
  selfSigned: {}
EOF

# install jaeger, elastic stack, and prometheus dependency
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

# install feelguuds_platform with tls enabled
helm upgrade --install feelguuds-platform ./charts/feelguuds_platform \
    --set image.repository=test/feelguuds_platform \
    --set image.tag=latest \
    --set tls.enabled=true \
    --set certificate.create=true \
    --namespace=default
