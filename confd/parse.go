package confd

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/docker/go-units"
	yaml "gopkg.in/yaml.v2"
)

// Percent %
type Percent int64

// Int64 return its value
func (p Percent) Int64() int64 {
	return int64(p)
}

// MarshalYAML implements the Marshaller interface.
func (p Percent) MarshalYAML() (interface{}, error) {
	return p, nil
}

// UnmarshalYAML implements the Unmarshaller interface.
func (p *Percent) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var percent string
	if err := unmarshal(&percent); err != nil {
		return err
	}

	percent = strings.TrimSuffix(percent, "%")
	val, err := strconv.ParseInt(percent, 10, 32)
	if err != nil {
		return err
	}
	*p = Percent(val * 1000)

	return nil
}

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
	CPULimit  Percent   `yaml:"cpu-limit"`
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
