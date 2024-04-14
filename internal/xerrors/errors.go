package xerrors

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInternal   = errors.New("internal error")
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
)

func Join(err1 error, err2 error) error {
	return fmt.Errorf("%w: %w", err1, err2)
}

func RenderHTML(w http.ResponseWriter, resource string, err error) {
	code, message := ErrorInfo(err)

	http.Error(w, fmt.Sprintf("%v: %v", resource, message), code)
}

func ErrorInfo(err error) (int, string) {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, ErrNotFound.Error()
	case errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest, ErrBadRequest.Error()
	}

	return http.StatusInternalServerError, ErrInternal.Error()
}
