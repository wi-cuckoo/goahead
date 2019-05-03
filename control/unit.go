package control

import (
	"errors"
	"os/exec"
	"strings"
	"syscall"
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
	if u.Dir != "" {
		cmd.Dir = u.Dir
	}
	cmd.Env = u.Envs

	return cmd, nil
}

// Kill the process of unit
func (u *Unit) Kill() error {
	if u.pid < 2 {
		return errors.New("invalid pid")
	}
	if err := syscall.Kill(u.pid, syscall.SIGINT); err != nil {
		return syscall.Kill(u.pid, syscall.SIGTERM)
	}
	return nil
}
