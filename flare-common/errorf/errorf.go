package errorf

import "fmt"

func ReadingFile(fileName string, err error) error {
	return fmt.Errorf("failed reading file %s with: %w", fileName, err)
}

func Unmarshal(fileName string, err error) error {
	return fmt.Errorf("failed marshaling file %s with: %w", fileName, err)
}
