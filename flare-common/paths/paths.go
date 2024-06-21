package paths

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// LocalToAbsolute returns the absolute path of the relativePath that is relative to the filepath of the caller.
//
// Example : if LocalToAbsolute(".../directory3/foo.txt") is called in file /User/dir1/dir2/foo.go,
// the return value is "/User/dir1/dir3/foo.go".
func LocalToAbsolute(relativePath string) (string, error) {

	_, b, _, ok := runtime.Caller(1)

	if !ok {
		return "", fmt.Errorf("LocalToAbsolute error for: %s", relativePath)
	}

	return filepath.Join(filepath.Dir(b), relativePath), nil

}
