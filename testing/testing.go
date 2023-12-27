package testing

import (
	"os"
	"path"
	"runtime"
)

// This init function changes the working directory to the root of the project for tests. This allows us to
// use relative paths in tests to get to testdata.
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	_ = os.Chdir(dir)
}
