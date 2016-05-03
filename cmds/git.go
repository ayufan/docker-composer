package cmds

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/ayufan/docker-composer/compose"
	"github.com/ayufan/docker-composer/helpers"
)

func runGitReceivePackCommand(c *cli.Context) error {
	app, err := compose.ExistingApplication(c.Args()...)
	if err != nil {
		logrus.Fatalln("Application:", err)
	}

	cmd := helpers.Git("receive-pack", app.Path())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func runGitUploadPackCommand(c *cli.Context) error {
	app, err := compose.ExistingApplication(c.Args()...)
	if err != nil {
		logrus.Fatalln("Application:", err)
	}

	cmd := helpers.Git("upload-pack", app.Path())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func init() {
	registerCommand(cli.Command{
		Name:      "git-receive-pack",
		Action:    runGitReceivePackCommand,
		Usage:     "receive changes from client",
		Category:  "git",
		ArgsUsage: "APP",
		HideHelp:  true,
	})

	registerCommand(cli.Command{
		Name:      "git-upload-pack",
		Action:    runGitUploadPackCommand,
		Usage:     "upload changes to client",
		Category:  "git",
		ArgsUsage: "APP",
		HideHelp:  true,
	})
}
