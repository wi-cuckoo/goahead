package control

// Resource limitation
type Resource struct {
	// Memory
	MemLimit int64 `yaml:"mem-limit"`
	MemResv  int64 `yaml:"mem-reservation"`
	// CPU
	CPUQuota int64 `yaml:"cpu-quota"`
}
