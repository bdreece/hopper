package utils

import "fmt"

func WrapError(outer, inner error) error {
	return fmt.Errorf("%w:\n%w", outer, inner)
}
