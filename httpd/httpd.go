package httpd

import (
	"net"
	"net/http"
	"net/url"

	"github.com/wi-cuckoo/goahead/control"
)

// Server handler
type Server struct {
	ln   net.Listener
	Ctrl control.Controller
}

// Start serve http handler on addr
func (s *Server) Start(addr string) error {
	_url, err := url.Parse(addr)
	if err != nil {
		return err
	}
	ln, err := net.Listen(_url.Scheme, _url.Path)
	if err != nil {
		return err
	}
	s.ln = ln

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello, man"))
	})

	go http.Serve(ln, nil)

	return nil
}

// Stop ...
func (s *Server) Stop() {
	if s.ln != nil {
		s.ln.Close()
	}
}
