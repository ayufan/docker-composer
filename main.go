package main

import (
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/ayufan/docker-composer/cmds"
	"github.com/ayufan/docker-composer/compose"
	"github.com/codegangsta/cli"
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
			Name:  "log-level, l",
			Value: "info",
			Usage: "Log level (options: debug, info, warn, error, fatal, panic)",
		},
		cli.StringFlag{
			Name:        "apps-dir",
			Value:       "/srv/apps",
			Usage:       "Directory where all the apps are stored",
			Destination: &compose.AppsDirectory,
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

	app.Commands = cmds.Commands

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
