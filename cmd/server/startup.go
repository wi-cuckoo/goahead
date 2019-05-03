package main

import (
	"errors"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/wi-cuckoo/goahead/control"
)

func run(c *cli.Context) error {
	ctrl, err := control.NewController("goahead")
	if err != nil {
		return err
	}
	defer ctrl.Destory()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit)
	<-exit

	logrus.Info("receive quit signal, exit ...")

	return errors.New("fuck you")
}
