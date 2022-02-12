package cmds

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/ayufan/docker-composer/compose"
	"github.com/ayufan/docker-composer/helpers"
)

func runEnvViewCommand(c *cli.Context) error {
	app, err := compose.ExistingApplication(c.Args()...)
	if err != nil {
		logrus.Fatalln("App:", err)
	}

	env, err := app.Env()
	if err != nil {
		logrus.Fatalln("App:", err)
	}

	println(env)
	return nil
}

func runEnvEditCommand(c *cli.Context) error {
	app, err := compose.ExistingApplication(c.Args()...)
	if err != nil {
		logrus.Fatalln("App:", err)
	}

	env, err := app.SupportsEnv()
	if err != nil {
		logrus.Fatalln("App:", err)
	}

	if !env {
		logrus.Fatalln("App does not support .env, likely as it is part of current repository")
	}

	for {
		err = helpers.EditFile(app.Path(".env"))
		if err != nil {
			logrus.Fatalln(err)
		}

		err = app.Compose("config", "-q")
		if err == nil {
			break
		}
		logrus.Errorln("Configuration is invalid. Press ENTER to edit again or CTRL-C to abort.")
		var input string
		fmt.Scanln(&input)
	}

	logrus.Infoln("Deploying...")
	err = app.Deploy()
	if err != nil {
		logrus.Fatalln("Deploy:", err)
	}

	logrus.Infoln("Tagging...")
	err = app.Tag()
	if err != nil {
		logrus.Fatalln("Tag:", err)
	}
	return nil
}

func init() {
	registerCommand(cli.Command{
		Name: "env",
		Subcommands: cli.Commands{
			cli.Command{
				Name:      "view",
				Action:    runEnvViewCommand,
				Usage:     "view environment configuration",
				Category:  "manage",
				ArgsUsage: "APP",
			},
			cli.Command{
				Name:      "edit",
				Action:    runEnvEditCommand,
				Usage:     "edit environment configuration",
				Category:  "manage",
				ArgsUsage: "APP",
			},
		},
		Usage:     "manage environment configuration",
		Category:  "manage",
		ArgsUsage: "APP",
	})
}
