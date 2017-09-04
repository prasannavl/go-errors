package goerror

type GoError interface {
	error
	Cause() error
}

func New(message string) GoError {
	return &GoErr{&message, nil}
}

func NewWithCause(message string, cause error) GoError {
	return &GoErr{&message, cause}
}

func From(err error) GoError {
	if err == nil {
		return nil
	}
	if gerr, ok := err.(GoError); ok {
		return gerr
	}
	return &GoErr{nil, err}
}

type GoErr struct {
	Msg   *string
	Inner error
}

func (e *GoErr) Error() string {
	if e.Msg != nil {
		return *e.Msg
	}
	cause := e.Cause()
	if cause != nil && cause != e {
		return cause.Error()
	}
	return "unknown error"
}

func (e *GoErr) Cause() error {
	return e.Inner
}