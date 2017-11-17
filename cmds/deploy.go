package cmds

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/ayufan/docker-composer/compose"
)

func deploySingleApplication(c *cli.Context, app *compose.App) (err error) {
	logrus.Println(app.Name, "...")
	if c.Bool("build") {
		err = app.Build()
	}
	if err == nil && c.Bool("pull") {
		err = app.Pull()
	}
	if err == nil {
		err = app.Deploy()
	}
	return
}

func runDeployCommand(c *cli.Context) error {
	var apps compose.AppList
	var err error

	if c.Bool("all") {
		apps, err = compose.Apps()
		apps = apps.OnlyEnabled()
	} else if c.NArg() != 0 {
		apps, err = compose.Apps(c.Args()...)
	} else {
		logrus.Fatalln("Specify at least one application")
	}
	if err != nil {
		logrus.Fatalln("Apps:", err)
	}

	for _, app := range apps {
		appErr := deploySingleApplication(c, app)
		if appErr != nil {
			logrus.Errorln(err)
			err = appErr
		}
	}

	if err != nil {
		logrus.Fatalln("Failed to deploy some of applications")
	}
	return nil
}

func init() {
	registerCommand(cli.Command{
		Name:      "deploy",
		Action:    runDeployCommand,
		Usage:     "deploy one or multiple applications",
		Category:  "global",
		ArgsUsage: "APP...",
		Flags: []cli.Flag{
			cli.BoolTFlag{
				Name:  "pull",
				Usage: "pull images",
			},
			cli.BoolTFlag{
				Name:  "build",
				Usage: "build images",
			},
			cli.BoolFlag{
				Name:  "all",
				Usage: "deploy all applications",
			},
		},
	})
}
