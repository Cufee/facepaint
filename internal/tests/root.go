package tests

import (
	"path/filepath"
	"runtime"
)

func Root() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	root, _ := filepath.Abs(filepath.Join(basepath, "../../"))
	return root
}
