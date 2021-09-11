#! /usr/bin/env sh

# add jetstack repository
kubectl apply -f https://raw.githubusercontent.com/pixie-labs/pixie/main/k8s/operator/crd/base/px.dev_viziers.yaml
kubectl apply -f https://raw.githubusercontent.com/pixie-labs/pixie/main/k8s/operator/helm/crds/olm_crd.yaml
helm repo add newrelic https://helm-charts.newrelic.com
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add jetstack https://charts.jetstack.io || true
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add kube-state-metrics https://kubernetes.github.io/kube-state-metrics
helm repo update

kubectl create namespace feelguuds-platform
kubectl create namespace newrelic

# install cert-manager
helm upgrade --install cert-manager jetstack/cert-manager \
    --set installCRDs=true \
    --namespace default

# wait for cert manager
kubectl cert-manager check api --wait=5m
kubectl rollout status deployment/cert-manager --timeout=5m
kubectl rollout status deployment/cert-manager-webhook --timeout=5m
kubectl rollout status deployment/cert-manager-cainjector --timeout=5m

# install self-signed certificate
cat << 'EOF' | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: self-signed
spec:
  selfSigned: {}
EOF

./install_charts

# install feelguuds_platform with tls enabled
helm upgrade --install feelguuds-platform ./charts/feelguuds_platform \
    --set image.repository=test/feelguuds_platform \
    --set image.tag=latest \
    --set tls.enabled=true \
    --set certificate.create=true \
    --namespace=default
