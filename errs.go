package abair

import (
	"fmt"
	"net/http"
)

// HttpErrorOptionFunc is an HTTP error option function.
type HttpErrorOptionFunc func(*HTTPError)

// HTTPError is an HTTP error.
type HTTPError struct {
	Code     int
	Message  any
	Internal error
}

// Error implements error.
func (h HTTPError) Error() string {
	if h.Internal == nil {
		return fmt.Sprintf("code: %d, message: %v", h.Code, h.Message)
	}
	return fmt.Sprintf("code: %d, message: %v, internal: %v", h.Code, h.Message, h.Internal)
}

// NewHTTPError creates a new HTTPError.
func NewHTTPError(code int, opts ...HttpErrorOptionFunc) *HTTPError {
	h := &HTTPError{
		Code:    code,
		Message: http.StatusText(code),
	}
	for _, opts := range opts {
		opts(h)
	}

	return h
}

// WithMessage sets the message.
func (h *HTTPError) WithMessage(message any) *HTTPError {
	h.Message = message
	return h
}

// WithInternal sets the internal error.
func (h *HTTPError) WithInternal(err error) *HTTPError {
	h.Internal = err
	return h
}
