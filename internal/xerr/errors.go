package xerr

import (
	"errors"
	"fmt"
)

var (
	ErrInternal = errors.New("internal error")
)

func Join(err1 error, err2 error) error {
	return fmt.Errorf("%w: %w", err1, err2)
}
