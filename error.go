package goerror

type Error interface {
	error
	Cause() error
	IsSource() bool
}

func New(message string) Error {
	return &Err{&message, nil}
}

func NewEx(message string, cause error) Error {
	return &Err{&message, cause}
}

func From(err error) Error {
	if err == nil {
		return nil
	}
	if gerr, ok := err.(Error); ok {
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
