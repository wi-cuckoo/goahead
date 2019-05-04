package main

import (
	"errors"
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
		EnvVar: "GOAHEAD_SERVER",
		Name:   "server",
		Usage:  "address to dail, eg: tcp://localhost:5555",
		Value:  "unix:///var/run/goahead.sock",
	},
}

var commands = []cli.Command{
	{
		Name:  "start",
		Usage: "start your program",
		Action: func(c *cli.Context) error {
			return run(c, "start")
		},
	},
	{
		Name:  "stop",
		Usage: "stop your program",
		Action: func(c *cli.Context) error {
			return run(c, "stop")
		},
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "goahead-cli"
	app.Usage = "control your application, like systemctl"
	app.Version = "unknown"
	app.Flags = flags
	app.Commands = commands
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func run(c *cli.Context, op string) error {
	program := c.Args().First()
	if program == "" {
		return errors.New("no program defined")
	}

	return nil
}

func init() {
	if uid := os.Geteuid(); uid != 0 {
		// panic("must be root")
	}
}
