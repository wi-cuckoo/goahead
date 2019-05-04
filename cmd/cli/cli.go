package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

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
	},
	{
		Name:  "stop",
		Usage: "stop your program",
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "goahead-cli"
	app.Usage = "control your application, like systemctl"
	app.Version = "unknown"
	app.Flags = flags
	app.Commands = commands
	app.Action = run
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func run(c *cli.Context) error {
	program := c.Args().First()
	if program == "" {
		return errors.New("no program defined")
	}

	_url, err := url.Parse(c.String("server"))
	if err != nil {
		return err
	}
	fmt.Println(c.String("server"), _url.Scheme)

	clnt := http.Client{
		Timeout: time.Second * 1,
		Transport: &http.Transport{
			Dial: func(network, addr string) (c net.Conn, err error) {
				fmt.Println(network, addr)
				return nil, nil
			},
			DisableKeepAlives: true,
		},
	}
	req := http.Request{
		Method: "PUT",
		URL:    _url,
	}
	resp, err := clnt.Do(&req)
	if err != nil {
		return err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	return nil
}

func init() {
	if uid := os.Geteuid(); uid != 0 {
		// panic("must be root")
	}
}
