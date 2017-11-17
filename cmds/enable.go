package cmds

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/ayufan/docker-composer/compose"
)

func runEnableCommand(c *cli.Context) error {
	app, err := compose.ExistingApplication(c.Args()...)
	if err != nil {
		logrus.Fatalln("App:", err)
	}

	err = app.Enable()
	if err != nil {
		logrus.Fatalln("Enable:", err)
	}

	err = app.Deploy()
	if err != nil {
		logrus.Fatalln("Deploy:", err)
	}
	return nil
}

func runDisableCommand(c *cli.Context) error {
	app, err := compose.ExistingApplication(c.Args()...)
	if err != nil {
		logrus.Fatalln("App:", err)
	}

	err = app.Disable()
	if err != nil {
		logrus.Fatalln("Disable:", err)
	}

	err = app.Compose("down")
	if err != nil {
		logrus.Fatalln(err)
	}
	return nil
}

func init() {
	registerCommand(cli.Command{
		Name:      "enable",
		Action:    runEnableCommand,
		Usage:     "enable application (previously disabled)",
		Category:  "manage",
		ArgsUsage: "APP",
	})
	registerCommand(cli.Command{
		Name:      "disable",
		Action:    runDisableCommand,
		Usage:     "disable application (previously enabled)",
		Category:  "manage",
		ArgsUsage: "APP",
	})
}
