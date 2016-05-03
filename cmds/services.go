package cmds

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/ayufan/docker-composer/compose"
)

func runServicesCommand(c *cli.Context) error {
	app, err := compose.ExistingApplication(c.Args()...)
	if err != nil {
		logrus.Fatalln("App:", err)
	}

	err = app.Compose("config", "--services")
	if err != nil {
		logrus.Fatalln("Compose:", err)
	}
	return nil
}

func init() {
	registerCommand(cli.Command{
		Name:      "services",
		Action:    runServicesCommand,
		Usage:     "list all services of application",
		Category:  "manage",
		ArgsUsage: "APP",
	})
}
