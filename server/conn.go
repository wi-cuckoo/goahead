package server

import (
	"encoding/json"
	"errors"
	"net"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wi-cuckoo/goahead"
	"github.com/wi-cuckoo/goahead/confd"
	"github.com/wi-cuckoo/goahead/control"
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
	var buf = make([]byte, 512)
	n, err := con.Read(buf)
	if err != nil {
		writelne(con, err)
		return
	}
	logrus.Info("recv from conn: ", string(buf[:n]))

	op := goahead.Operation{}
	if err := json.Unmarshal(buf[:n], &op); err != nil {
		writelne(con, err)
		return
	}
	switch op.Command {
	case "start":
		// start a program
		s.startProgram(con, op.Program)
	case "stop":
		// stop a program
		s.stopProgram(con, op.Program)
	case "status":
		s.statusProgram(con, op.Program)
	default:
		// unknown
		writelne(con, errors.New("invalid command"))
	}
	logrus.Info("done: ", op.String())
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
		writelne(con, err)
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
