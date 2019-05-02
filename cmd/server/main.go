package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var flags = []cli.Flag{
	cli.BoolFlag{
		EnvVar: "GOAHEAD_DEBUG",
		Name:   "debug",
		Usage:  "enable debug mode, default false",
	},
	cli.StringFlag{
		EnvVar: "GOAHEAD_CONF",
		Name:   "config, c",
		Usage:  "config file to startup daemon",
		Value:  "/etc/goahead.conf",
	},
}

var commands = []cli.Command{}

func main() {
	app := cli.NewApp()
	app.Name = "goahead"
	app.Usage = "control your application, like systemd"
	app.Version = "unknown"
	app.Flags = flags
	app.Action = run
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
