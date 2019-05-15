package control

import (
	"bytes"
	"strings"
	"text/template"
)

// Status for an unit
type Status struct {
	Uptime string
	PID    int
	CPU    string
	Mem    string
}

var statsTmpl = `
uptime		: {{.Uptime}}
pid		: {{.PID}}
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

	return strings.TrimSpace(buf.String())
}
