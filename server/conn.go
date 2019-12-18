package server

import (
	"net"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wi-cuckoo/goahead/confd"
	"github.com/wi-cuckoo/goahead/control"
	"github.com/wi-cuckoo/goahead/pb"
	"github.com/wi-cuckoo/goahead/util"
)

const sock = "/var/run/goahead.sock"

func init() {
	os.Remove(sock)
}

// SocketServer to listen
type SocketServer struct {
	ln   net.Listener
	Conf *confd.Store
	Ctrl control.Controller
}

// Start serve http handler on addr
func (s *SocketServer) Start() error {
	ln, err := net.Listen("unix", sock)
	if err != nil {
		return err
	}
	s.ln = ln

	go func() {
		for {
			con, err := ln.Accept()
			if err != nil {
				break
			}
			go s.handleConn(con)
		}
	}()

	return nil
}

func (s *SocketServer) handleConn(con net.Conn) {
	defer con.Close()
	decoder := util.NewDecoder(con, 1<<9)
	in, err := decoder.DecodeInstruct()
	if err != nil {
		writelne(con, err)
		return
	}
	logrus.Info("recv instruct:", in.String())
	switch in.Op {
	case pb.Op_START:
		// start a program
		s.startProgram(con, in.App)
	case pb.Op_STOP:
		// stop a program
		s.stopProgram(con, in.App)
	case pb.Op_STATUS:
		s.statusProgram(con, in.App)
	default:
		// unknown
		writeln(con, "invalid command")
	}
	logrus.Infof("%s %s done", in.Op, in.App)
}

// Stop ...
func (s *SocketServer) Stop() {
	if s.ln != nil {
		s.ln.Close()
	}
}

func (s *SocketServer) startProgram(con net.Conn, name string) {
	cfg, err := s.Conf.GetConfig(name)
	if err != nil {
		writeln(con, err.Error())
		return
	}
	unit := control.Unit{
		Name:  name,
		Owner: cfg.Owner,
		Desc:  cfg.Desc,
		Dir:   cfg.Directory,
		Envs:  cfg.Envs,
		Cmd:   cfg.Command,
		Res: &control.Resource{
			CPUQuota: cfg.CPULimit.Int64(),
			MemLimit: cfg.MemLimit.Int64(),
		},
	}

	errCh := make(chan error)
	go func() {
		if err := s.Ctrl.Start(&unit); err != nil {
			errCh <- err
		}
	}()

	writeln(con, "goahead starting "+name)
	select {
	case err := <-errCh:
		writelne(con, err)
	case <-time.After(time.Second * 3):
		writeln(con, "goahead started "+name)
	}
}

func (s *SocketServer) stopProgram(con net.Conn, name string) {
	if err := s.Ctrl.Stop(name); err != nil {
		writelne(con, err)
		return
	}
	writeln(con, "goahead stopped "+name)
}

func (s *SocketServer) statusProgram(con net.Conn, name string) {
	status, err := s.Ctrl.Status(name)
	if err != nil {
		writelne(con, err)
		return
	}
	writeln(con, status.String())
}
