package vyos

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrMethodNotSupported = errors.New("method not supported")
	ErrContextNil         = errors.New("context must be non-nil")
	ErrInterfaceNil       = errors.New("can not unmarshal into nil interface")
	ErrEmptyPath          = errors.New("path cannot be empty")
	ErrMustLoadFromFile   = errors.New("file must not be empty or nil")
)

type HTTPError struct {
	code int
	body string
}

func (h *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %s: %s", http.StatusText(h.code), h.body)
}
