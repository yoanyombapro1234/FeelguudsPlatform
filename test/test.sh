#1 /usr/bin/env sh

set -e

# wait for feelguuds platform
kubectl rollout status deployment/feelguuds_platform --timeout=3m

# test feelguuds platform
helm test feelguuds_platform
