package errors

// wraps the original error with a new error
// that has a Public() string method to return
// a message acceptable for display in a toast
func Public(err error, msg string) error {
	return publicError{err, msg}
}

type publicError struct {
	err error
	msg string
}

func (e publicError) Error() string {
	return e.err.Error()
}

func (e publicError) Public() string {
	return e.msg
}

func (e publicError) Unwrap() error {
	return e.err
}
