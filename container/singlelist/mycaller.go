package list

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func MyCallerFileLine() string {

	fpcs := make([]uintptr, 1)

	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return "n/a"
	}

	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "n/a"
	}

	filename, line := fun.FileLine(fpcs[0] - 1)
	filename = filepath.Base(filename)
	return fmt.Sprintf("%v:%v", filename, line)
}
