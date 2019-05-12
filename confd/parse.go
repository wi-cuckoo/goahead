package confd

import (
	"io/ioutil"

	"github.com/docker/go-units"
	yaml "gopkg.in/yaml.v2"
)

// SpaceSize for disk, memory
type SpaceSize int64

// Int64 return its value
func (s SpaceSize) Int64() int64 {
	return int64(s)
}

// MarshalYAML implements the Marshaller interface.
func (s SpaceSize) MarshalYAML() (interface{}, error) {
	return s, nil
}

// UnmarshalYAML implements the Unmarshaller interface.
func (s *SpaceSize) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var size string
	err := unmarshal(&size)
	if err != nil {
		return err
	}
	val, err := units.RAMInBytes(size)
	if err != nil {
		return err
	}
	*s = SpaceSize(val)

	return nil
}

// ProgramConfig ...
type ProgramConfig struct {
	Name      string    `yaml:"name"`
	Owner     string    `yaml:"owner"`
	Desc      string    `yaml:"description"`
	Directory string    `yaml:"directory"`
	Command   string    `yaml:"command"`
	Envs      []string  `yaml:"environments"`
	CPUQuota  int64     `yaml:"cpu-quota"`
	MemLimit  SpaceSize `yaml:"mem-limit"`
	modify    int64     `yaml:"_"` // modify time
}

func parseYamlFile(filename string) (*ProgramConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := ProgramConfig{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
