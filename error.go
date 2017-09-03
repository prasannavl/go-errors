package goerror

type GoError interface {
	error
	Cause() error
	IsSource() bool
}

func New(message string) GoError {
	return &Err{&message, nil}
}

func NewWithCause(message string, cause error) GoError {
	return &Err{&message, cause}
}

func From(err error) GoError {
	if err == nil {
		return nil
	}
	if gerr, ok := err.(GoError); ok {
		return gerr
	}
	return &Err{nil, err}
}

type Err struct {
	Msg   *string
	Inner error
}

func (e *Err) Error() string {
	if e.Msg != nil {
		return *e.Msg
	}
	cause := e.Cause()
	if cause != nil && cause != e {
		return cause.Error()
	}
	return "unknown error"
}

func (e *Err) IsSource() bool {
	return e.Msg != nil
}

func (e *Err) Cause() error {
	return e.Inner
}
