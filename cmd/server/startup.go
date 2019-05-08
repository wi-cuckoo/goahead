package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/wi-cuckoo/goahead/confd"
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

	conf, err := confd.NewStore(c.String("dir"))
	if err != nil {
		return err
	}

	s := server.SocketServer{
		Conf: conf,
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
