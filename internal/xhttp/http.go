package xhttp

import (
	"errors"
	"net/http"

	"github.com/tombuente/apex/internal/xerrors"
)

func ErrorToStatusCode(err error) uint {
	switch {
	case errors.Is(err, xerrors.ErrNotFound):
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}
