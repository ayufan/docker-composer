package cmds

import (
	"github.com/codegangsta/cli"

	"github.com/ayufan/docker-composer/helpers"
)

func runCleanupCommand(c *cli.Context) error {
	if c.BoolT("containers") {
		helpers.System("docker ps -f status=exited -f status=dead -a --format '{{.ID}}\t{{.Labels}}' | grep -v com.docker.compose | awk '{print $1}' | xargs -r docker rm -f")
	}
	if c.BoolT("run-containers") {
		helpers.System("docker ps -f label=com.docker.compose.oneoff=True -aq | xargs -r docker rm -f")
	}
	if c.BoolT("volumes") {
		helpers.System("docker volume ls -qf dangling=true | xargs -r docker volume rm")
	}
	if c.BoolT("images") {
		helpers.System("docker images -qf dangling=true | xargs -r docker rmi")
	}
	if c.BoolT("unused-images") {
		helpers.System("docker images -q | xargs -r docker rmi")
	}
	return nil
}

func init() {
	registerCommand(cli.Command{
		Name:     "cleanup",
		Action:   runCleanupCommand,
		Usage:    "cleanup containers, docker images and volumes",
		Category: "maintenance",
		Flags: []cli.Flag{
			cli.BoolTFlag{
				Name:  "containers",
				Usage: "cleanup exited or dead containers",
			},
			cli.BoolTFlag{
				Name:  "run-containers",
				Usage: "cleanup docker-compose run containers",
			},
			cli.BoolTFlag{
				Name:  "volumes",
				Usage: "cleanup unused volumes",
			},
			cli.BoolTFlag{
				Name:  "images",
				Usage: "cleanup dangling images",
			},
			cli.BoolTFlag{
				Name:  "unused-images",
				Usage: "cleanup unused images",
			},
		},
	})
}
