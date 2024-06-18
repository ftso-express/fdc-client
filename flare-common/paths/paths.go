package paths

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func LocalToAbsolute(relativePath string) (string, error) {

	_, b, _, ok := runtime.Caller(1)

	fmt.Println(b)

	if !ok {
		return "", fmt.Errorf("LocalToAbsolutePath error for: %s", relativePath)
	}

	return filepath.Join(filepath.Dir(b), relativePath), nil

}
