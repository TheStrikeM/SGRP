package e

import "fmt"

func Wrap(operation string, err error) error {
	return fmt.Errorf("%s: %w", operation, err)
}
