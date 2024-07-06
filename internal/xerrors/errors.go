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
	message, code := HttpInfo(err)

	http.Error(w, fmt.Sprintf("%v: %v", resource, message), code)
}

func HttpInfo(err error) (string, int) {
	switch {
	case errors.Is(err, ErrNotFound):
		return ErrNotFound.Error(), http.StatusNotFound
	case errors.Is(err, ErrBadRequest):
		return ErrBadRequest.Error(), http.StatusBadRequest
	}

	return ErrInternal.Error(), http.StatusInternalServerError
}
