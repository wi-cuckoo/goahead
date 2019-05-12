package control

import (
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// Resource limitation
type Resource struct {
	// Memory
	MemLimit int64 `yaml:"mem-limit"`
	MemResv  int64 `yaml:"mem-reservation"`
	// CPU
	CPUQuota int64 `yaml:"cpu-quota"`
}

// Convert2Specs ...
func (r *Resource) Convert2Specs() *specs.LinuxResources {
	return &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Quota: &r.CPUQuota,
		},
		Memory: &specs.LinuxMemory{
			Limit:       &r.MemLimit,
			Reservation: &r.MemResv,
		},
	}
}
