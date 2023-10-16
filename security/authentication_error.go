package security

type (
	WrongUserNamePasswordError struct {
		Err error
	}
	NotPromissionError struct {
		Err error
	}
	MethodNotAllowedError struct {
		Err error
	}
	AccountLockedError struct {
		Err error
	}
	AccountExpiredError struct {
		Err error
	}
)

func (e WrongUserNamePasswordError) Error() string {
	return e.Err.Error()
}

func (e NotPromissionError) Error() string {
	return e.Err.Error()
}

func (e MethodNotAllowedError) Error() string {
	return e.Err.Error()
}

func (e AccountLockedError) Error() string {
	return e.Err.Error()
}

func (e AccountExpiredError) Error() string {
	return e.Err.Error()
}
