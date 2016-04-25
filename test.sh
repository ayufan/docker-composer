#!/bin/bash

set -e
docker build -t composer . &>/dev/null
docker run -i -t -v /var/run/docker.sock:/var/run/docker.sock --rm composer "$@"
