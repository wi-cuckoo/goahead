package control

import (
	"errors"
	"syscall"
)

// Kill the process of unit
func (u *Unit) Kill() error {
	if u.pid < 2 {
		return errors.New("invalid pid")
	}
	if err := syscall.Kill(u.pid, syscall.SIGINT); err != nil {
		return syscall.Kill(u.pid, syscall.SIGKILL)
	}
	return nil
}
