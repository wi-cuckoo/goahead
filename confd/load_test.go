package confd

import "testing"

func TestNewConfd(t *testing.T) {
	t.Log(NewConfd("/etc/goahead.d"))
}
