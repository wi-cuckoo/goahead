package confd

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// ProgramConfig ...
type ProgramConfig struct {
	Name       string   `yaml:"name"`
	Maintainer string   `yaml:"maintainer"`
	Directory  string   `yaml:"directory"`
	Command    string   `yaml:"command"`
	Envs       []string `yaml:"environments"`
	CPUQuota   int64    `yaml:"cpu-quota"`
	MemLimit   int64    `yaml:"mem-limit"`
	modify     int64    `yaml:"_"` // modify time
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
