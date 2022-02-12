package cmds

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/ryanuber/columnize"

	"github.com/ayufan/docker-composer/compose"
)

func runLsCommand(c *cli.Context) error {
	detail := c.Bool("detail")

	apps, err := compose.Apps(c.Args()...)
	if err != nil {
		logrus.Fatalln(err)
	}

	if detail {
		lines := []string{"APP | STATUS | REF"}
		for _, app := range apps {
			status, _ := app.Status()
			reference, _ := app.Reference()
			lines = append(lines, fmt.Sprintf("%s | %s | %s", app.Name, status, reference))
		}
		println(columnize.SimpleFormat(lines))
	} else {
		for _, app := range apps {
			println(app.Name)
		}
	}
	return nil
}

func init() {
	command := cli.Command{
		Name:      "ls",
		Action:    runLsCommand,
		Usage:     "show all applications",
		ArgsUsage: "APP...",
		Category:  "global",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "detail",
				Usage: "show detailed status of application",
			},
		},
	}
	registerCommand(command)
}
