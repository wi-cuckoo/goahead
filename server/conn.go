package server

import (
	"net"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/wi-cuckoo/goahead/control"
)

const sock = "/tmp/goahead.sock"

func init() {
	os.Remove(sock)
}

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
	logrus.Info("recv con: ", string(buf[:n]))
	con.Write([]byte("done, man!"))
}

// Stop ...
func (s *SocketServer) Stop() {
	if s.ln != nil {
		s.ln.Close()
	}
}
