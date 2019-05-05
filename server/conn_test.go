package server

import (
	"io"
	"net"
	"testing"
)

func TestSockServer(t *testing.T) {
	ss := &SocketServer{}
	if err := ss.Start(); err != nil {
		t.Fatal(err)
	}
	defer ss.Stop()

	con, _ := net.Dial("unix", sock)
	con.Write([]byte("start yourprogram"))
	for {
		buf := make([]byte, 512)
		n, err := con.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error(err)
			break
		}
		t.Log(string(buf[:n]))
	}
}
