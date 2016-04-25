package cmds

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/Sirupsen/logrus"

	"github.com/ayufan/docker-composer/compose"
)

func runInitCommand(c *cli.Context) {
	apps, err := compose.Application(c.Args()...)
	if err != nil {
		logrus.Fatalln(err)
	}

	err = apps.Init()
	if err != nil {
		logrus.Fatalln(err)
	}
}

func runReleaseCommand(c *cli.Context) {
	app, err := compose.ExistingApplication(c.Args()...)
	if err != nil {
		logrus.Fatalln(err)
	}

	args := []string{"--remove-orphans"}
	if c.Bool("volumes") {
		args = append(args, "--volumes")
	}
	if c.String("rmi") != "" {
		args = append(args, "--rmi", c.String("rmi"))
	}

	err = app.Compose("down", args...)
	if err != nil {
		if c.Bool("force") {
			logrus.Errorln(err)
		} else {
			logrus.Fatalln(err)
		}
	}

	err = os.RemoveAll(app.Path())
	if err != nil {
		if c.Bool("force") {
			logrus.Errorln(err)
		} else {
			logrus.Fatalln(err)
		}
	}
}

func init() {
	registerCommand(cli.Command{
		Name:      "init",
		Action:    runInitCommand,
		Usage:     "initialize a new application",
		Category:  "manage",
		ArgsUsage: "APP",
	})
	registerCommand(cli.Command{
		Name:      "release",
		Action:    runReleaseCommand,
		Usage:     "release a new application",
		Category:  "manage",
		ArgsUsage: "APP",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name: "force",
				Usage: "Unconditionally remove",
			},
			cli.BoolFlag{
				Name: "volumes",
				Usage: "Remove data volumes",
			},
			cli.StringFlag{
				Name: "rmi",
				Usage: "Remove images, type may be one of: 'all' to remove all images, or 'local' to remove only images that don't have an custom name set by the `image` field",
			},
		},
	})
}
