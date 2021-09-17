#!/usr/bin/env bash

# installing postgres helm chart under a given release name
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

kubectl create namespace newrelic

# install cert-manager
helm upgrade --install cert-manager jetstack/cert-manager \
    --set installCRDs=true \
    --namespace default

# wait for cert manager
kubectl rollout status deployment/cert-manager --timeout=3m
kubectl rollout status deployment/cert-manager-webhook --timeout=3m
kubectl rollout status deployment/cert-manager-cainjector --timeout=3m

# install self-signed certificate
cat << 'EOF' | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: self-signed
spec:
  selfSigned: {}
EOF

helm upgrade --install prometheus prometheus-community/prometheus

# create a namespace for new relic bundle and deploy app
helm upgrade --install newrelic-bundle newrelic/nri-bundle \
 --set global.licenseKey="$NEW_RELIC_LICENSE" \
 --set global.cluster=development \
 --namespace=newrelic \
 --set newrelic-infrastructure.privileged=true \
 --set ksm.enabled=true \
 --set prometheus.enabled=true \
 --set kubeEvents.enabled=true \
 --set logging.enabled=true \
 --set newrelic-pixie.enabled=true \
 --set newrelic-pixie.apiKey="$NEW_RELIC_PIXIE_DEPLOY" \
 --set pixie-chart.enabled=true \
 --set pixie-chart.deployKey="$NEW_RELIC_PIXIE_APIKEY" \
 --set pixie-chart.clusterName=development

# we build the feelguuds docker image and send it to minikube registry which will be pulled by the by helm
# during deployment
# link: https://medium.com/swlh/how-to-run-locally-built-docker-images-in-kubernetes-b28fbc32cc1d
make mkd_push_image

# install feelguuds platform and all dependent services
helm upgrade --install feelguuds-platform ./charts/feelguuds-platform-dev
