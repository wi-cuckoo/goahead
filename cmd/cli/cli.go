package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/urfave/cli"
	"github.com/wi-cuckoo/goahead/pb"
	"github.com/wi-cuckoo/goahead/util"
)

var flags = []cli.Flag{
	cli.StringFlag{
		EnvVar: "GOAHEAD_SOCK",
		Name:   "sock",
		Usage:  "sock address to dail",
		Value:  "/var/run/goahead.sock",
	},
}

var commands = []cli.Command{
	{
		Name:  "start",
		Usage: "start your program",
		Action: func(c *cli.Context) error {
			return run(c, pb.Op_START)
		},
	},
	{
		Name:  "stop",
		Usage: "stop your program",
		Action: func(c *cli.Context) error {
			return run(c, pb.Op_STOP)
		},
	},
	{
		Name:  "status",
		Usage: "status your program",
		Action: func(c *cli.Context) error {
			return run(c, pb.Op_STATUS)
		},
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
	app.Name = "goahead-cli"
	app.Usage = "control your application, like supervisor"
	app.Version = version
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("version=%s revision=%s\n", c.App.Version, revision)
	}
	app.Flags = flags
	app.Commands = commands
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func run(c *cli.Context, op pb.Op) error {
	program := c.Args().First()
	if program == "" {
		return errors.New("no program defined")
	}

	con, err := net.DialTimeout("unix", c.GlobalString("sock"), time.Second*2)
	if err != nil {
		return err
	}
	defer con.Close()

	encoder := util.NewEncoder(con, 1<<9)
	if err := encoder.EncodeInstruct(&pb.Instruct{
		Op:  op,
		App: program,
	}); err != nil {
		return err
	}
	if err := encoder.Flush(); err != nil {
		return err
	}
	con.SetReadDeadline(time.Now().Add(time.Second * 10))

	_, err = io.Copy(os.Stdout, con)
	return err
}

func init() {
	if uid := os.Geteuid(); uid != 0 {
		panic("must be root")
	}
}
