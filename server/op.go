package server

// Operation ...
type Operation struct {
	Command string `json:"command"`
	Program string `json:"program"`
}

func (o *Operation) String() string {
	return o.Command + " " + o.Program
}
