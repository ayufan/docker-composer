package cmds

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/ayufan/docker-composer/compose"
)

func runUpdateHooksCommand(c *cli.Context) error {
	apps, err := compose.Apps(c.Args()...)
	if err != nil {
		logrus.Fatalln(err)
	}

	for _, app := range apps {
		logrus.Infoln(app.Name, "...")
		err := app.UpdateHooks()
		if err != nil {
			logrus.Fatalln("Hooks:", err)
		}

		err = app.UpdateConfig()
		if err != nil {
			logrus.Fatalln("Config:", err)
		}
	}
	return nil
}

func init() {
	registerCommand(cli.Command{
		Name:      "update-hooks",
		Action:    runUpdateHooksCommand,
		Usage:     "update hooks of applications",
		Category:  "maintenance",
		ArgsUsage: "APP...",
	})
}
