package compose

const pushToCheckout = `#!/usr/bin/env bash

set -eo pipefail

trap 'git reset --hard' EXIT
echo "Applying new changes..."
git reset --hard >/dev/null
git update-index -q --refresh
git read-tree -u -m HEAD "$1"
echo "Deploying application..."
docker-compose up -d --build --remove-orphans
trap - EXIT
git tag "$1" -f latest
`
