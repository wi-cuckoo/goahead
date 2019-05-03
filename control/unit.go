package control

import (
	"errors"
	"os/exec"
	"strings"
)

// Unit for a process
type Unit struct {
	Name  string
	Desc  string
	Owner string

	Dir  string
	Cmd  string
	Envs []string

	Res Resource

	pid int
}

// Command return shell command
func (u *Unit) Command() (*exec.Cmd, error) {
	args := strings.Split(u.Cmd)
	if len(args) < 1 {
		return nil, errors.New("invalid command")
	}

	cmd := exec.Command(args[0], args[1:])
}
