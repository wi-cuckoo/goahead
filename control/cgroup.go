package control

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/containerd/cgroups"
	"github.com/docker/go-units"
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

	ctrl, err := c.root.New(u.Name, u.Res.Convert2Specs())
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
	// 保存状态
	u.up = time.Now()
	u.pid = pid
	u.ctrl = ctrl
	c.Store(u.Name, u)
	defer c.Delete(u.Name)

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
	el, ok := c.Load(id)
	if !ok {
		return nil, fmt.Errorf("%s not found", id)
	}
	u, _ := el.(*Unit)

	metrics, err := u.ctrl.Stat(cgroups.IgnoreNotExist)
	if err != nil {
		return nil, err
	}

	cpu, mem := metrics.CPU, metrics.Memory
	cpuUsage := float64(cpu.Usage.Total) / float64(cpu.Throttling.Periods) / 1000000
	return &Status{
		Uptime: units.HumanDuration(time.Since(u.up)),
		PID:    u.pid,
		CPU:    fmt.Sprintf("%.2f%%", cpuUsage),
		Mem:    units.HumanSize(float64(mem.Usage.Usage)),
	}, nil
}

func (c *cgCtrl) Destory() error {
	c.Range(func(key, val interface{}) bool {
		if u, ok := val.(*Unit); ok {
			u.Kill()
		}
		return true
	})
	return c.root.Delete()
}

// NewController return a cgCtrl instance
func NewController(root string) (Controller, error) {
	path := "/" + root
	ctrl, err := cgroups.New(cgroups.V1, cgroups.StaticPath(path), &specs.LinuxResources{})
	if err != nil {
		ctrl, err = cgroups.Load(cgroups.V1, cgroups.StaticPath(path))
		if err != nil {
			return nil, err
		}
	}
	cc := &cgCtrl{
		root: ctrl,
	}
	return cc, nil
}
