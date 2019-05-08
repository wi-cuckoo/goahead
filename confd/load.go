package confd

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Store relationship between file and its content
type Store struct {
	mux sync.RWMutex
	dir string
	rel map[string]*ProgramConfig
}

// NewStore ...
func NewStore(dir string) (*Store, error) {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
	}
	if err != nil {
		return nil, err
	}
	rel := make(map[string]*ProgramConfig)
	for _, f := range filesWithSuffixInDir(dir, ".yaml") {
		cfg, err := parseYamlFile(filepath.Join(dir, f.Name()))
		if err != nil {
			continue
		}
		cfg.modify = f.ModTime().Unix()
		rel[f.Name()] = cfg
	}

	return &Store{
		dir: dir,
		rel: rel,
	}, nil
}

// GetConfig by program name
func (s *Store) GetConfig(program string) (*ProgramConfig, error) {
	filename := filepath.Join(s.dir, program+".yaml")
	f, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	s.mux.RLock()
	old, ok := s.rel[f.Name()]
	s.mux.RUnlock()

	if ok && old.modify == f.ModTime().Unix() {
		return old, nil
	}

	cfg, err := parseYamlFile(filename)
	if err != nil {
		return nil, err
	}
	cfg.modify = f.ModTime().Unix()

	s.mux.Lock()
	s.rel[f.Name()] = cfg
	s.mux.Unlock()

	return cfg, nil
}

func filesWithSuffixInDir(dir, suffix string) []os.FileInfo {
	files := make([]os.FileInfo, 0, 100)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(filepath.Base(path), suffix) {
			return nil
		}

		files = append(files, info)

		return nil
	})
	return files
}
