package cmds

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/ayufan/docker-composer/compose"
	"github.com/ayufan/docker-composer/helpers"
)

func gitEditor() (editor string, err error) {
	editor = os.Getenv("GIT_EDITOR")
	term := os.Getenv("TERM")
	isDumb := term == "" || term == "dumb"

	if editor == "" && !isDumb {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" && isDumb {
		return "", errors.New("No GIT_EDITOR defined")
	}
	if editor == "" {
		editor = "vi"
	}
	return
}

func editFile(name string) (err error) {
	editor, err := gitEditor()
	if err != nil {
		return
	}

	logrus.Infoln("Editing", filepath.Base(name), "...")
	cmd := helpers.Command(editor, name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func runEditCommand(c *cli.Context) error {
	app, err := compose.ExistingApplication(c.Args()...)
	if err != nil {
		logrus.Fatalln("App:", err)
	}

	for {
		err = editFile(app.Path("docker-compose.yml"))
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
