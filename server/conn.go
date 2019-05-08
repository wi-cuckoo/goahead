package server

import (
	"encoding/json"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wi-cuckoo/goahead/confd"
	"github.com/wi-cuckoo/goahead/control"
)

const sock = "/var/run/goahead.sock"

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
		con.Write([]byte(err.Error()))
		return
	}
	logrus.Info("recv from conn: ", string(buf[:n]))

	op := Operation{}
	if err := json.Unmarshal(buf[:n], &op); err != nil {
		con.Write([]byte("invalid command"))
		return
	}
	switch op.Command {
	case "start":
		// start a program
		s.startProgram(con, op.Program)
	case "stop":
		// stop a program
		con.Write([]byte("stopping " + op.Program))
		<-time.After(time.Second * 3)
		con.Write([]byte("stopped " + op.Program))
	default:
		// unknown
		con.Write([]byte("unknown command"))
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
		con.Write([]byte(err.Error()))
		return
	}
	unit := control.Unit{
		Name:  name,
		Owner: cfg.Owner,
		Desc:  cfg.Desc,
		Dir:   cfg.Directory,
		Envs:  cfg.Envs,
		Cmd:   cfg.Command,
		Res: control.Resource{
			CPUQuota: cfg.CPUQuota,
			MemLimit: cfg.MemLimit,
		},
	}

	errCh := make(chan error)
	go func() {
		if err := s.Ctrl.Start(&unit); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-errCh:
		con.Write([]byte(err.Error()))
	case <-time.After(time.Second * 3):
		con.Write([]byte("started " + name))
	}
}
