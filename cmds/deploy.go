package cmds

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/ayufan/docker-composer/compose"
)

func runDeployCommand(c *cli.Context) {
	var apps []*compose.App
	var err error

	if c.Bool("all") {
		apps, err = compose.Apps()
	} else if c.NArg() != 0 {
		apps, err = compose.Apps(c.Args()...)
	} else {
		logrus.Fatalln("Specify at least one application")
	}
	if err != nil {
		logrus.Fatalln("Apps:", err)
	}

	for _, app := range apps {
		logrus.Println(app.Name, "...")
		appErr := app.Deploy()
		if appErr != nil {
			err = appErr
		}
		if appErr != nil {
			logrus.Errorln(appErr)
		}
	}

	if err != nil {
		logrus.Fatalln("Failed to deploy some of applications")
	}
}

func init() {
	registerCommand(cli.Command{
		Name:      "deploy",
		Action:    runDeployCommand,
		Usage:     "deploy one or multiple applications",
		Category:  "global",
		ArgsUsage: "APP...",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "all",
				Usage: "deploy all applications",
			},
		},
	})
}
