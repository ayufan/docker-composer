package main

import (
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/mattn/go-shellwords"

	"github.com/ayufan/docker-composer/cmds"
	"github.com/ayufan/docker-composer/compose"
)

func main() {
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "a Docker Composer Service"
	app.Author = "Kamil Trzciński"
	app.Email = "ayufan@ayufan.eu"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug mode",
			EnvVar: "DEBUG",
		},
		cli.StringFlag{
			Name:   "log-level, l",
			Value:  "info",
			Usage:  "Log level (options: debug, info, warn, error, fatal, panic)",
			EnvVar: "LOG_LEVEL",
		},
		cli.StringFlag{
			Name:        "apps-dir",
			Value:       "/srv/apps",
			Usage:       "Directory where all the apps are stored",
			Destination: &compose.AppsDirectory,
			EnvVar:      "APPS_DIR",
		},
		cli.StringFlag{
			Name:  "c",
			Usage: "Custom command to execute",
		},
	}

	// logs
	app.Before = func(c *cli.Context) error {
		logrus.SetOutput(os.Stderr)
		level, err := logrus.ParseLevel(c.String("log-level"))
		if err != nil {
			logrus.Fatalf(err.Error())
		}
		logrus.SetLevel(level)
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors: true,
		})

		// If a log level wasn't specified and we are running in debug mode,
		// enforce log-level=debug.
		if !c.IsSet("log-level") && !c.IsSet("l") && c.Bool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}

	defaultAction := app.Action
	app.Action = func(c *cli.Context) error {
		if command := c.String("c"); command != "" {
			args, err := shellwords.Parse(command)
			if err != nil {
				logrus.Fatalln(err)
			}

			args = append([]string{os.Args[0]}, args...)
			return app.Run(args)
		}

		return cli.HandleAction(defaultAction, c)
	}

	app.Commands = cmds.Commands
	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatalln(err)
	}
}
