#1 /usr/bin/env sh

set -e

# wait for podinfo
kubectl rollout status deployment/feelguuds_platform --timeout=3m

# test podinfo
helm test feelguuds_platform
