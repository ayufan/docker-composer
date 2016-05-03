package cmds

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/ayufan/docker-composer/compose"
)

func runSwitchCommand(c *cli.Context) error {
	if c.NArg() != 2 {
		logrus.Fatalln("Specify APP and REF")
	}

	app, err := compose.ExistingApplication(c.Args()[0])
	if err != nil {
		logrus.Fatalln("App:", err)
	}

	err = app.Reset(c.Args()[1])
	if err != nil {
		logrus.Fatalln("Git:", err)
	}

	err = app.Deploy()
	if err != nil {
		app.Revert()
		logrus.Fatalln("Compose:", err)
	}

	app.Tag()
	return nil
}

func init() {
	registerCommand(cli.Command{
		Name:      "switch",
		Action:    runSwitchCommand,
		Usage:     "switch application to different branch",
		Category:  "manage",
		ArgsUsage: "APP REF",
	})
}
