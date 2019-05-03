package main

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/wi-cuckoo/goahead/control"
)

// ROOTCGROUP define cgroup hierarchy root
const ROOTCGROUP = "goahead"

func run(c *cli.Context) error {
	ctrl, err := control.NewController(ROOTCGROUP)
	if err != nil {
		return err
	}
	defer ctrl.Destory()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	<-exit

	logrus.Info("receive quit signal, exit ...")

	return errors.New("fuck you")
}

