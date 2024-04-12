package xerrors

import (
	"errors"
	"fmt"
)

var (
	ErrInternal = errors.New("internal error")
	ErrNotFound = errors.New("not found")
)

func Join(err1 error, err2 error) error {
	return fmt.Errorf("%w: %w", err1, err2)
}
