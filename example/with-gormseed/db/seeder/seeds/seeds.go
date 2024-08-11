package seeds

import (
	"path/filepath"
	"runtime"
)

type Seeds struct{}

func (*Seeds) Path() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	return basepath
}
