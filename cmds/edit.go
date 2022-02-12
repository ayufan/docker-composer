package cmds

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/ayufan/docker-composer/compose"
	"github.com/ayufan/docker-composer/helpers"
)

func runEditCommand(c *cli.Context) error {
	app, err := compose.ExistingApplication(c.Args()...)
	if err != nil {
		logrus.Fatalln("App:", err)
	}

	for {
		err = helpers.EditFile(app.Path("docker-compose.yml"))
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

	logrus.Infoln("Addding docker-compose.yml...")
	err = app.Git("add", "docker-compose.yml")
	if err != nil {
		logrus.Fatalln("Add:", err)
	}

	logrus.Infoln("Checking for changes...")
	changed, err := app.Changed()
	if err != nil {
		logrus.Fatalln("Status:", err)
	} else if !changed {
		logrus.Infoln("Nothing changed.")
		return nil
	}

	logrus.Infoln("Deploying...")
	err = app.Deploy()
	if err != nil {
		logrus.Fatalln("Deploy:", err)
	}

	logrus.Infoln("Committing...")
	err = app.Commit("manually editted")
	if err != nil {
		logrus.Fatalln("Commit:", err)
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
		Name:      "edit",
		Action:    runEditCommand,
		Usage:     "edit configuration of application",
		Category:  "manage",
		ArgsUsage: "APP",
	})
}
