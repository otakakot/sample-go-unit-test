package errors

import "errors"

type NotFoundError struct {
	err error
}

func NewNotFoundError(err error) NotFoundError {
	return NotFoundError{
		err: err,
	}
}

func (nfe NotFoundError) Error() string {
	if nfe.err != nil {
		return nfe.err.Error()
	}

	return "not found"
}

func AsNotFoundError(err error) bool {
	var target NotFoundError

	return errors.As(err, &target)
}
