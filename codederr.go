package goerror

type CodedError interface {
	Error
	Code() int
}

func NewCoded(code int, message string) CodedError {
	return &CodedErr{Err{&message, nil}, code}
}

func NewCodedEx(code int, message string, cause error) CodedError {
	return &CodedErr{Err{&message, cause}, code}
}

func CodedFrom(code int, err error) CodedError {
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
