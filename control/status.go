package control

import (
	"bytes"
	"text/template"
	"time"
)

// Status for an unit
type Status struct {
	Uptime time.Duration
	PID    int
	CPU    string
	Mem    string
}

var statsTmpl = `
uptime		: {{.Uptime}}
pid			: {{.PID}}
cpu usage	: {{.CPU}}
mem usage	: {{.Mem}}
`

func (s Status) String() string {
	tmpl, err := template.New("status").Parse(statsTmpl)
	if err != nil {
		return err.Error()
	}
	var buf = new(bytes.Buffer)
	tmpl.Execute(buf, s)

	return buf.String()
}
