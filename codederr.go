package goerror

type CodedError interface {
	GoError
	Code() int
}

func NewCoded(code int, message string) CodedError {
	return &CodedErr{Err{&message, nil}, code}
}

func NewCodedWithCause(code int, message string, cause error) CodedError {
	return &CodedErr{Err{&message, cause}, code}
}

func CodedFrom(err error, code int) CodedError {
	if err == nil {
		return nil
	}
	if gerr, ok := err.(CodedError); ok && code == gerr.Code() {
		return gerr
	}
	return &CodedErr{Err{nil, err}, code}
}

type CodedErr struct {
	Err
	ErrCode int
}

func (h *CodedErr) Code() int {
	return h.ErrCode
}
