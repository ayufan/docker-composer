#!/bin/sh
#
# Run docker-composer in a container

set -e

export DOCKER_IMAGE="ayufan/docker-composer:latest"
DOCKER_RUN_OPTIONS=""
DOCKER_ADDR=""
VOLUMES="-v /srv/apps:/srv/apps"

# Setup options for connecting to docker host
if [ -z "$DOCKER_HOST" ]; then
    DOCKER_HOST="/var/run/docker.sock"
fi
if [ -S "$DOCKER_HOST" ]; then
    DOCKER_ADDR="-v $DOCKER_HOST:$DOCKER_HOST -e DOCKER_HOST"
else
    DOCKER_ADDR="-e DOCKER_HOST -e DOCKER_TLS_VERIFY -e DOCKER_CERT_PATH"
fi

# Only allocate tty if we detect one
DOCKER_RUN_OPTIONS="-i"
if [ -t 1 ]; then
    DOCKER_RUN_OPTIONS="$DOCKER_RUN_OPTIONS -t"
fi

exec docker run --rm $DOCKER_RUN_OPTIONS -e DOCKER_IMAGE $DOCKER_ADDR $VOLUMES $DOCKER_IMAGE "$@"
