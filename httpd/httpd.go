package httpd

import (
	"net"
	"net/http"
	"net/url"

	"github.com/wi-cuckoo/goahead/control"
)

// Server handler
type Server struct {
	Ctrl control.Controller
}

// NewServer serve http handler on addr
func NewServer(addr string) (*Server, error) {
	_url, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	ln, err := net.Listen(_url.Scheme, _url.Path)
	if err != nil {
		return nil, err
	}

	http.Serve(ln, nil)

	return nil, nil
}
