package confd

import (
	"io/ioutil"
	"sync"

	yaml "gopkg.in/yaml.v2"
)

// Store relationship between file and its content
type Store struct {
	mux sync.Mutex
	rel map[string]*ProgramConfig
}

// ProgramConfig ...
type ProgramConfig struct {
	Name       string   `yaml:"name"`
	Maintainer string   `yaml:"maintainer"`
	Directory  string   `yaml:"directory"`
	Command    string   `yaml:"command"`
	Envs       []string `yaml:"environments"`
	CPUQuota   int64    `yaml:"cpu-quota"`
	MemLimit   int64    `yaml:"mem-limit"`
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
