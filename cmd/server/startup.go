package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/wi-cuckoo/goahead/control"
	"github.com/wi-cuckoo/goahead/server"
)

// ROOTCGROUP define cgroup hierarchy root
const ROOTCGROUP = "goahead"

func run(c *cli.Context) error {
	ctrl, err := control.NewController(ROOTCGROUP)
	if err != nil {
		return err
	}
	defer ctrl.Destory()

	s := server.SocketServer{
		Ctrl: ctrl,
	}
	if err := s.Start(); err != nil {
		return err
	}
	defer s.Stop()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	<-exit

	logrus.Info("receive quit signal, exit ...")

	return nil
}
