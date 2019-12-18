package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/wi-cuckoo/goahead"
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
		Value:  "/etc/goahead.d",
	},
}

var (
	// Revision ...
	revision string
	// Version ...
	version string
)

func main() {
	app := cli.NewApp()
	app.Name = "goahead"
	app.Usage = "control your application, like supervisor"
	app.Version = version
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("version=%s revision=%s\n", c.App.Version, revision)
	}
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
	fmt.Fprintf(os.Stdout, "%s\t%s\n", goahead.Banner, version)
}
