package cmds

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

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

func runGitPushToCheckoutCommand(c *cli.Context) (err error) {
	if len(c.Args()) < 2 {
		logrus.Fatalln("Usage:", os.Args[0])
	}

	app, err := compose.ExistingApplication(c.Args()[0])
	if err != nil {
		logrus.Fatalln("Application:", err)
	}

	revision := c.Args()[1]
	if revision == "" {
		logrus.Fatalln("--revision is required")
	}

	logrus.Infoln("Applying new changes...")

	defer func() {
		if err != nil {
			logrus.Warningln("Restoring repository state...")
			app.Git("reset", "--hard")
		}
	}()

	err = app.Git("update-index", "-q", "--refresh")
	if err != nil {
		logrus.Fatalln("Update-index:", err)
	}

	err = app.Git("read-tree", "-u", "--reset", revision)
	if err != nil {
		logrus.Fatalln("Read-tree:", err)
	}

	logrus.Infoln("Deploying application...")
	err = app.Build()
	if err != nil {
		logrus.Fatalln("Build:", err)
	}

	err = app.Deploy()
	if err != nil {
		logrus.Fatalln("Deploy:", err)
	}

	logrus.Infoln("Tagging release...")
	err = app.Tag(revision)
	if err != nil {
		logrus.Fatalln("Tag:", err)
	}

	return nil
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

	registerCommand(cli.Command{
		Name:      "git-push-to-checkout",
		Action:    runGitPushToCheckoutCommand,
		Usage:     "push to checkout deployment hook",
		Category:  "git",
		ArgsUsage: "APP ...",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "revision",
				Usage: "revision to deploy",
			},
		},
		HideHelp: true,
	})
}
