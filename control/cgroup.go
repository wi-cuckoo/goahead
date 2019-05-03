package control

import (
	"errors"
	"fmt"
	"sync"

	"github.com/containerd/cgroups"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

type cgCtrl struct {
	sync.Map // name -> unit
	root     cgroups.Cgroup
}

func (c *cgCtrl) Start(u *Unit) error {
	if _, ok := c.Load(u.Name); ok {
		return errors.New("exist, use reload maybe")
	}

	ctrl, err := c.root.New(u.Name, &specs.LinuxResources{})
	if err != nil {
		return err
	}
	defer ctrl.Delete()

	// start command
	cmd, err := u.Command()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	pid := cmd.Process.Pid
	if err := ctrl.Add(cgroups.Process{Pid: pid}); err != nil {
		return err
	}
	u.pid = pid
	c.Store(u.Name, u)

	return cmd.Wait()
}

func (c *cgCtrl) Stop(id string) error {
	el, ok := c.Load(id)
	if !ok {
		return fmt.Errorf("%s not found", id)
	}
	u, _ := el.(*Unit)

	return u.Kill()
}

func (c *cgCtrl) Reload(id string) error {
	return nil
}

func (c *cgCtrl) Status(id string) (*Status, error) {
	return nil, nil
}

// NewController return a cgCtrl instance
func NewController(root string) (Controller, error) {
	path := "/" + root
	ctrl, err := cgroups.New(cgroups.V1, cgroups.StaticPath(path), &specs.LinuxResources{})
	if err != nil {
		ctrl, err = cgroups.Load(cgroups.V1, path)
		if err != nil {
			return nil, err
		}
	}
	cc := &cgCtrl{
		root: ctrl,
	}
	return cc, nil
}
