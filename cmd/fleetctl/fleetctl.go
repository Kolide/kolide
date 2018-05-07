package main

import (
	"math/rand"
	"time"

	"github.com/kolide/kit/version"
	"github.com/urfave/cli"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	app := cli.NewApp()
	app.Name = "fleetctl"
	app.Usage = "The CLI for operating Kolide Fleet"
	app.Version = version.Version().Version
	cli.VersionPrinter = func(c *cli.Context) {
		version.PrintFull()
	}

	app.Commands = []cli.Command{
		applyCommand(),
		setupCommand(),
		loginCommand(),
		logoutCommand(),
		cli.Command{
			Name:  "get",
			Usage: "Get/list resources",
			Subcommands: []cli.Command{
				getQueriesCommand(),
			},
		},
		cli.Command{
			Name:  "config",
			Usage: "Modify how and which Fleet server to connect to",
			Subcommands: []cli.Command{
				configSetCommand(),
				configGetCommand(),
			},
		},
	}

	app.RunAndExitOnError()
}
