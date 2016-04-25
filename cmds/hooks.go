package cmds

import (
	"github.com/codegangsta/cli"
	"github.com/Sirupsen/logrus"

	"github.com/ayufan/docker-composer/compose"
)

func runUpdateHooksCommand(c *cli.Context) {
	apps, err := compose.Apps(c.Args()...)
	if err != nil {
		logrus.Fatalln(err)
	}

	for _, app := range apps {
		logrus.Infoln(app.Name, "...")
		err := app.UpdateHooks()
		if err != nil {
			logrus.Fatalln(err)
		}
	}
}

func init() {
	registerCommand(cli.Command{
		Name:      "update-hooks",
		Action:    runUpdateHooksCommand,
		Usage:     "update hooks of applications",
		Category:  "git",
		ArgsUsage: "APP...",
	})
}