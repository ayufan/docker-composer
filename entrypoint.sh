#!/bin/bash

set -eo pipefail

cd "$APPS_DIR/"

case "$1" in)
	build|config|create|down|exec|kill|logs|pause|port|ps|pull|restart|rm|run|scale|start|stop|unpause|up)
		if [[ $# -le 2 ]]; then
			echo "Please specify APP" 1>&2
			exit
		fi
		CMD="$1"
		APP="$2"
		shift 2
		cd "$APP/"
		git init &>/dev/null
		docker-compose "$CMD" "$@"
		;;

	edit)
		APP="$2"
		shift 2
		cd "$APP/"
		git init &>/dev/null
		nano docker-compose.yml
		git add docker-compose.yml
		if git commit -am "configuration updated"; then
			docker-compose up --build --remove-orphans -d
		fi
		;;

	events|version)
		docker-compose "$@"
		;;

	help)
		echo Commands:
		cat <<EOF
  build APP          Build or rebuild services
  config APP         Validate and view the compose file
  create APP         Create services
  down APP           Stop and remove containers, networks, images, and volumes
  events             Receive real time events from containers
  exec APP           Execute a command in a running container
  help               Get help on a command
  kill APP           Kill containers
  logs APP           View output from containers
  pause APP          Pause services
  port APP           Print the public port for a port binding
  ps  APP            List containers
  pull APP           Pulls service images
  restart APP        Restart services
  rm APP             Remove stopped containers
  run APP            Run a one-off command
  scale APP          Set number of containers for a service
  start APP          Start services
  stop APP           Stop services
  unpause APP        Unpause services
  up APP             Create and start containers
  version            Show the Docker-Compose version information
EOF
		;;

	apps)
		
esac
