package util

import (
	"bufio"
	"errors"
	"io"
	"strconv"

	"github.com/gogo/protobuf/proto"
	"github.com/wi-cuckoo/goahead/pb"
)

// Encoder ...
type Encoder struct {
	bw *bufio.Writer
}

// NewEncoder ..
func NewEncoder(w io.Writer, size int) *Encoder {
	return &Encoder{bufio.NewWriterSize(w, size)}
}

// EncodeInstruct ...
func (e *Encoder) EncodeInstruct(in *pb.Instruct) error {
	b, err := proto.Marshal(in)
	if err != nil {
		return err
	}
	head := make([]byte, 0, 1<<3)
	head = append(head, pre)
	n := len(b)
	head = append(head, []byte(strconv.Itoa(n))...)
	head = append(head, cr)
	if _, err := e.bw.Write(head); err != nil {
		return err
	}
	_, err = e.bw.Write(b)
	return err
}

// Flush call underlayer bufio flush
func (e *Encoder) Flush() error {
	return e.bw.Flush()
}

const (
	pre byte = '$'
	cr  byte = '\r'
)

// Decoder ...
type Decoder struct {
	br *bufio.Reader
}

// NewDecoder ...
func NewDecoder(r io.Reader, size int) *Decoder {
	return &Decoder{bufio.NewReaderSize(r, size)}
}

// DecodeInstruct ...
func (d *Decoder) DecodeInstruct() (*pb.Instruct, error) {
	if _, err := d.readpre(); err != nil {
		return nil, err
	}
	n, err := d.readint()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, n)
	if _, err := io.ReadFull(d.br, buf); err != nil {
		return nil, err
	}
	in := new(pb.Instruct)
	err = proto.Unmarshal(buf, in)
	return in, err
}

func (d *Decoder) readpre() (byte, error) {
	b, err := d.br.ReadByte()
	if b != pre {
		return b, errors.New("invlid prefix byte")
	}
	return b, err
}

func (d *Decoder) readint() (int, error) {
	bs, err := d.br.ReadBytes(cr)
	if err != nil {
		return 0, err
	}
	if len(bs) < 2 {
		return 0, errors.New("invaid instruct length encoding")
	}
	return strconv.Atoi(string(bs[:len(bs)-1]))
}
