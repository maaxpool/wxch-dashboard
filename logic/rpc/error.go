package rpc

import "errors"

type HTTPError struct {
	error

	errorCode uint32
}

func warpHTTPError(err error) *HTTPError {
	return warpHTTPErrorWithCode(500, err)
}

func warpHTTPErrorWithCode(errorCode uint32, err error) *HTTPError {
	return &HTTPError{
		error:     err,
		errorCode: errorCode,
	}
}

func NewHttpError(errorCode uint32, err string) *HTTPError {
	return &HTTPError{
		error:     errors.New(err),
		errorCode: errorCode,
	}
}

func NewHttpRedirect(path string) *HTTPError {
	return &HTTPError{
		error:     errors.New(path),
		errorCode: 302,
	}
}
