package confd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// NewConfd ...
func NewConfd(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
	}
	if err != nil {
		return err
	}
	files := make([]os.FileInfo, 0, 100)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(filepath.Base(path), ".yaml") {
			return nil
		}

		files = append(files, info)

		return nil
	})
	for _, f := range files {
		fmt.Println(f.Name())
	}
	return nil
}
