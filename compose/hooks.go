package compose

const pushToCheckout = `#!/usr/bin/env bash

set -eo pipefail

cd ..

trap 'git reset --hard' EXIT
echo "Applying new changes..."
git update-index -q --refresh
git read-tree -u --reset "$1"
echo "Deploying application..."
docker-compose up -d --build --remove-orphans
trap - EXIT
`

const postReceive = `#!/usr/bin/env bash

set -eo pipefail
git tag "$2" -f latest
`
