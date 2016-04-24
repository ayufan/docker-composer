#!/bin/bash

set -eo pipefail

cd "$APPS_DIR/"

require_app() {
	if [[ $# -lt 2 ]]; then
		echo "Please specify APP" 1>&2
		exit
	fi
}

containers() {
	docker-compose ps -q 2>/dev/null
}

health() {
	STATUS=( $(docker inspect --format '{{.State.Running}}' "$@" 2>/dev/null | sort -u) )
	if [[ -n "${STATUS[1]}" ]]; then
		echo partial
	elif [[ "${STATUS[0]}" == "true" ]]; then
		echo running
	elif [[ "${STATUS[0]}" == "true" ]]; then
		echo not running
	else
		echo unknown
	fi
}

case "$1" in
	all)
		ls -1
		;;

	status)
		ls -1 | while read APP; do
			(
				cd $APP
				echo -e "$APP\t$(health $(containers))"
			)
		done
		;;

	all-detail)
		ls -1 | while read APP; do
			(
				cd $APP/
				docker-compose ps
			)
		done
		;;

	deploy-all)
		ls -1 | xargs -0 $0 deploy 
		;;

	create)
		require_app "$@"
		git init "$2/"
		cd "$2/"
		shift 2
		;;

	destroy)
		require_app "$@"
		APP="$2"
		shift 2
		(
			cd "$APP/"
			docker-compose down --remove-orphans --rmi local --volumes
		)
		rm -r "$APP/"
		;;

	deploy)
		require_app "$@"
		cd "$2/"
		shift 2
		docker-compose up --build --remove-orphans -d
		git tag -f latest
		;;

	edit)
		require_app "$@"
		cd "$2/"
		shift 2
		nano "docker-compose.yml"
		echo Validating configuration...
		docker-compose config &>/dev/null
		git init
		git add "docker-compose.yml"
		if [[ -z $(git status -s) ]]; then
			echo "Nothing changed."
			exit 0
		fi

		echo Updating configuration...
		docker-compose up --build --remove-orphans -d
		git commit -am "configuration updated"
		git tag -f latest
		;;

	revert)
		require_app "$@"
		cd "$2/"
		shift 2
		git reset --hard latest
		docker-compose up --build --remove-orphans -d
		;;

	stats)
		require_app "$@"
		APP="$2"
		shift 2
		docker stats --all "$@" $(docker inspect --format '{{.Name}}' $(containers "$APP"))
		;;

	services)
		require_app "$@"
		cd "$2/"
		shift 2
		docker-compose config --services "$@"
		;;

	git-*)
		require_app "$@"
		cd "$2/"
		git init
		GIT="$1"
		shift 2
		"$GIT" "$@"
		;;		

	build|config|down|exec|kill|logs|pause|port|ps|pull|restart|rm|run|scale|start|stop|unpause|up)
		require_app "$@"
		CMD="$1"
		APP="$2"
		shift 2
		cd "$APP/"
		git init &>/dev/null
		docker-compose "$CMD" "$@"
		;;

	events|version)
		docker-compose "$@"
		;;

	help)
		cat <<EOF
Usage:
  $0 command [options]

Composer Commands:
  all                List applications
  status             Show status of all applications
  create APP         Initialize application
  destroy APP        Destroy application
  edit APP           Edit configuration
  revert APP         Revert configuration
  deploy APP         Deploy application
  stats APP          Show application statistics
  services APP       Show application services
  deploy-all         Re-deploy all applications

Docker Compose Commands:
  build APP          Build or rebuild services
  config APP         Validate and view the compose file
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

	"")
		$0 help
		;;

	*)
		echo "Command not found: $@" 1>&2
		exit 1
		;;	
esac
