package goerror

type CodedError interface {
	Error
	Code() int
	Fatal() bool
}

func NewCoded(code int, message string, isFatal bool) CodedError {
	return &CodedErr{Err{&message, nil}, code, isFatal}
}

func NewCodedEx(code int, message string, cause error, isFatal bool) CodedError {
	return &CodedErr{Err{&message, cause}, code, isFatal}
}

func CodedFrom(code int, err error, isFatal bool) CodedError {
	if err == nil {
		return nil
	}
	if gerr, ok := err.(CodedError); ok && code == gerr.Code() {
		return gerr
	}
	return &CodedErr{Err{nil, err}, code, isFatal}
}

type CodedErr struct {
	Err
	ErrCode int
	IsFatal bool
}

func (h *CodedErr) Code() int {
	return h.ErrCode
}

func (h *CodedErr) Fatal() bool {
	return h.IsFatal
}
