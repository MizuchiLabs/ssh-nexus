// Package test contains test utilities
package test

import (
	"os"
	"path"
	"runtime"
)

// fixes the current working directory for tests
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
