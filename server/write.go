package server

import (
	"errors"
	"net"
)

func writeln(con net.Conn, s string) (int, error) {
	buf := []byte(s)
	if buf[len(buf)-1] != '\n' {
		buf = append(buf, '\n')
	}
	return con.Write(buf)
}

func writelne(con net.Conn, err error) (int, error) {
	if err == nil {
		err = errors.New("empty error")
	}
	buf := []byte(err.Error())
	if buf[len(buf)-1] != '\n' {
		buf = append(buf, '\n')
	}
	return con.Write(buf)
}
