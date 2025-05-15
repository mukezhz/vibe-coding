package utils

import (
	"os"
	"path/filepath"
)

func ChDir(path ...string) {
	if len(path) == 0 {
		path = []string{"../.."}
	}
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	newPath := filepath.Join(cwd, path[0])
	err = os.Chdir(newPath)
	if err != nil {
		panic(err)
	}
}
