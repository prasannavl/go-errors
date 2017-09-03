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
	if err == nil {
		return nil
	}
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
