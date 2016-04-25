package cmds

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/ayufan/docker-composer/compose"
)

var composeAppCommands map[string]string = map[string]string{
	"build":   "Build or rebuild services",
	"config":  "Validate and view the compose file",
	"create":  "Create services",
	"down":    "Stop and remove containers, networks, images, and volumes",
	"events":  "Receive real time events from containers",
	"exec":    "Execute a command in a running container",
	"kill":    "Kill containers",
	"logs":    "View output from containers",
	"pause":   "Pause services",
	"port":    "Print the public port for a port binding",
	"ps":      "List containers",
	"pull":    "Pulls service images",
	"restart": "Restart services",
	"rm":      "Remove stopped containers",
	"run":     "Run a one-off command",
	"scale":   "Set number of containers for a service",
	"start":   "Start services",
	"stop":    "Stop services",
	"unpause": "Unpause services",
	"up":      "Create and start containers",
}

func runComposeCommand(c *cli.Context) {
	if c.NArg() < 1 {
		logrus.Fatalln("Missing application name")
	}

	app, err := compose.ExistingApplication(c.Args()[0])
	if err != nil {
		logrus.Fatalln("App:", err)
	}

	err = app.Compose(c.Command.Name, c.Args()[1:]...)
	if err != nil {
		logrus.Fatalln("Compose:", err)
	}
}

func init() {
	for commandName, commandHelp := range composeAppCommands {
		command := cli.Command{
			Name:            commandName,
			Action:          runComposeCommand,
			SkipFlagParsing: true,
			Usage:           commandHelp,
			ArgsUsage:       "APP",
			Category:        "compose",
		}
		registerCommand(command)
	}
}
