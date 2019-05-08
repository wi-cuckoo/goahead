package confd

import "testing"

func TestNewConfd(t *testing.T) {
	t.Log(NewStore("/etc/goahead.d"))
}
