package goerror

type CodedError interface {
	GoError
	Code() int
}

func NewCoded(code int, message string) CodedError {
	return &CodedErr{GoErr{&message, nil}, code}
}

func NewCodedWithCause(code int, message string, cause error) CodedError {
	return &CodedErr{GoErr{&message, cause}, code}
}

func CodedFrom(err error, code int) CodedError {
	// If nil, just return back a nil so that
	// it can easily be composed without having
	// to check for nil first
	if err == nil {
		return nil
	}
	// Avoid a potentially wasteful allocation,
	// if it's already the same error with the
	// same error props
	if gerr, ok := err.(CodedError); ok && code == gerr.Code() {
		return gerr
	}
	return &CodedErr{GoErr{nil, err}, code}
}

type CodedErr struct {
	GoErr
	ErrCode int
}

func (h *CodedErr) Code() int {
	return h.ErrCode
}
