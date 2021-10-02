package publicerror

import (
	"errors"
	"net/http"
)

/*
Error handling that accounts for the difference between public and private
error messages.

This is a transparent error wrapping library using a custom error type which
allows us to inject a "public" error message at any point in the error cause
chain without breaking support for errors.Is() or errors.Unwrap().

The most common use-case is server responses for failed operations in which
you want to return a different message to the client than the one you log
with sensitive application details.

Example:

	err = publicerror.Wrap(err, "Opps! Not Found!", http.StatusNotFound)

	...

	// Log all errors and nested errors
	log.Println(err)

	// Show the client a "safe" error or http.StatusInternalServerError
	http.Error(w, publicerror.Message(err), publicerror.StatusCode(err))

*/

type Error struct {
	Err        error
	Message    string
	StatusCode int
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) Error() string {
	return e.Err.Error()
}

// Wrap error with public-safe message and status code (if err is not null)
func Wrap(err error, message string, code int) error {
	if err == nil {
		return nil
	}

	return Error{
		Err:        err,
		Message:    message,
		StatusCode: code,
	}
}

// Find and return first publicerror.Error in error chain
func Find(err error) *Error {
	if err == nil {
		return nil
	} else if e, ok := err.(Error); ok {
		return &e
	} else if ok && e.Err != nil {
		return Find(e.Err)
	} else if err2 := errors.Unwrap(err); err2 != nil {
		return Find(err2)
	}
	return nil
}

// StatusCode of first publicerror.Error in error chain, if available.
// Otherwise returns http.StatusInternalServerError
func StatusCode(err error) int {
	if e := Find(err); e != nil {
		return e.StatusCode
	}
	return http.StatusInternalServerError
}

// Message of first publicerror.Error in error chain, if available.
// Otherwise returns http.StatusText(http.StatusInternalServerError)
func Message(err error) string {
	if e := Find(err); e != nil {
		return e.Message
	}
	return http.StatusText(http.StatusInternalServerError)
}
