package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var flags = []cli.Flag{
	cli.BoolFlag{
		EnvVar: "GOAHEAD_DEBUG",
		Name:   "debug",
		Usage:  "enable debug mode, default false",
	},
	cli.StringFlag{
		EnvVar: "GOAHEAD_CONFDIR",
		Name:   "dir, d",
		Usage:  "dir to load config file of subprocess",
		Value:  "/etc/goaheadconf.d",
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "goahead"
	app.Usage = "control your application, like systemd"
	app.Version = "unknown"
	app.Flags = flags
	app.Action = run
	if err := app.Run(os.Args); err != nil {
		logrus.Error("app.Run err: ", err)
	}
}

func init() {
	if uid := os.Geteuid(); uid != 0 {
		panic("must be root")
	}
}
