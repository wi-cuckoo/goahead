package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/wi-cuckoo/goahead/server"

	"github.com/urfave/cli"
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

func run(c *cli.Context, cmd string) error {
	program := c.Args().First()
	if program == "" {
		return errors.New("no program defined")
	}

	con, err := net.DialTimeout("unix", c.GlobalString("sock"), time.Second*2)
	if err != nil {
		return err
	}

	buf, _ := json.Marshal(server.Operation{cmd, program})
	if _, err := con.Write(buf); err != nil {
		return err
	}

	con.SetReadDeadline(time.Now().Add(time.Second * 10))

	br := bufio.NewReader(con)
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Error("con.Read: ", err)
			break
		}
		fmt.Fprintln(os.Stdout, line)
	}

	return nil
}

func init() {
	if uid := os.Geteuid(); uid != 0 {
		panic("must be root")
	}
}
