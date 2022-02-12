package cmds

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/ayufan/docker-composer/compose"
	"github.com/ayufan/docker-composer/helpers"
)

func runStatsCommand(c *cli.Context) error {
	apps, err := compose.Apps(c.Args()...)
	if err != nil {
		logrus.Fatalln(err)
	}

	var names []string
	for _, app := range apps {
		appNames, err := app.ContainerNames()
		if err != nil {
			continue
		}
		names = append(names, appNames...)
	}
	if len(names) == 0 {
		logrus.Fatalln("No containers.")
	}

	cmd := helpers.Command("docker", "stats")
	cmd.Args = append(cmd.Args, names...)
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		logrus.Fatalln("Run:", err)
	}
	return nil
}

func init() {
	registerCommand(cli.Command{
		Name:      "stats",
		Action:    runStatsCommand,
		Usage:     "show statistics of all applications",
		Category:  "global",
		ArgsUsage: "APP...",
	})
}
