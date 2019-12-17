package util_test

import (
	"bytes"
	"testing"

	"github.com/wi-cuckoo/goahead/pb"
	"github.com/wi-cuckoo/goahead/util"
)

func TestEncode(t *testing.T) {
	in := &pb.Instruct{
		Op:  pb.Op_START,
		App: "my-app",
	}
	rw := new(bytes.Buffer)
	encoder := util.NewEncoder(rw, 128)
	if err := encoder.EncodeInstruct(in); err != nil {
		t.Fatal(err)
	}
	encoder.Flush()
	decoder := util.NewDecoder(rw, 128)
	iin, err := decoder.DecodeInstruct()
	if err != nil {
		t.Fatal(err)
	}
	if iin.Op != in.Op || iin.App != in.App {
		t.Error("decoded instruct not match")
	}
}
