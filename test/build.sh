#! /usr/bin/env sh

set -e

# build the docker file
GIT_COMMIT=$(git rev-list -1 HEAD) && \
DOCKER_BUILDKIT=1 docker build --tag test/feelguuds_platform --build-arg "REVISION=${GIT_COMMIT}" .
