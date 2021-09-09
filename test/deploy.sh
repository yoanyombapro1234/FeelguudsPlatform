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

./install_helm_charts

# install feelguuds_platform with tls enabled
helm upgrade --install feelguuds-platform ./charts/feelguuds_platform \
    --set image.repository=test/feelguuds_platform \
    --set image.tag=latest \
    --set tls.enabled=true \
    --set certificate.create=true \
    --namespace=default
