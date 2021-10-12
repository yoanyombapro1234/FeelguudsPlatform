#1 /usr/bin/env sh

set -e

# wait for feelguuds platform
kubectl rollout status deployment/feelguuds-platform --timeout=3m

# test feelguuds platform
helm test feelguuds-platform
