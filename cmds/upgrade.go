package cmds

import (
	"os"

	"github.com/codegangsta/cli"

	"github.com/ayufan/docker-composer/helpers"
)

var dockerImage = os.Getenv("DOCKER_IMAGE")

func runUpgradeCommand(c *cli.Context) error {
	cmd := helpers.Docker("docker", "pull", dockerImage)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func init() {
	if dockerImage != "" {
		registerCommand(cli.Command{
			Name:     "self-upgrade",
			Action:   runUpgradeCommand,
			Usage:    "upgrade docker-composer version",
			Category: "maintenance",
		})
	}
}
