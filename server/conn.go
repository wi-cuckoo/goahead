package server

import (
	"encoding/json"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/wi-cuckoo/goahead/control"
)

const sock = "/var/run/goahead.sock"

// SocketServer to listen
type SocketServer struct {
	ln   net.Listener
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
		con.Write([]byte("start done"))
	case "stop":
		// stop a program
		con.Write([]byte("stop done"))
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
