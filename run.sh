#!/bin/sh
#
# Run docker-composer in a container

set -e

SCRIPT_PATH="$(readlink -f "$0")"

export DOCKER_IMAGE="ayufan/docker-composer:${DOCKER_COMPOSER_TAG-latest}"
DOCKER_RUN_OPTIONS="-e GIT_EDITOR -e EDITOR -e VISUAL -e TERM"
DOCKER_ADDR=""
VOLUMES="-v /srv/apps:/srv/apps"

# Install user if missing
if ! id -u compose &>/dev/null; then
    sudo useradd -m -G docker -s "$SCRIPT_PATH" compose
    sudo install -d -m 700 ~compose/.ssh -o compose -g compose
    if [ -e ~/.ssh/authorized_keys ]; then
      sudo cp ~/.ssh/authorized_keys ~compose/.ssh/authorized_keys
      sudo chown compose:compose ~compose/.ssh/authorized_keys
      sudo chmod 600 ~compose/.ssh/authorized_keys
    fi
fi

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
DOCKER_RUN_OPTIONS="$DOCKER_RUN_OPTIONS -i"
if [ -t 1 ]; then
    DOCKER_RUN_OPTIONS="$DOCKER_RUN_OPTIONS -t"
fi

exec docker run --rm $DOCKER_RUN_OPTIONS -e DOCKER_IMAGE $DOCKER_ADDR $VOLUMES $DOCKER_IMAGE "$@"
