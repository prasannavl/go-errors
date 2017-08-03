package httperror

import http "net/http"
import goerror "github.com/prasannavl/goerror"

type HttpError interface {
	goerror.CodedError
	Stop() bool
}

func New(code int, message string, stop bool) HttpError {
	return &HttpErr{goerror.CodedErr{goerror.Err{&message, nil}, ErrorCode(code)}, stop}
}

func NewEx(code int, message string, cause error, stop bool) HttpError {
	return &HttpErr{goerror.CodedErr{goerror.Err{&message, cause}, ErrorCode(code)}, stop}
}

func From(code int, err error, stop bool) HttpError {
	if err == nil {
		return nil
	}
	code = ErrorCode(code)
	if gerr, ok := err.(HttpError); ok &&
		code == gerr.Code() && stop == gerr.Stop() {
		return gerr
	}
	return &HttpErr{goerror.CodedErr{goerror.Err{nil, err}, code}, stop}
}

type HttpErr struct {
	goerror.CodedErr
	ShouldStop bool
}

func (h *HttpErr) Stop() bool {
	return h.ShouldStop
}

func StatusCode(code int) int {
	if IsStatusCode(code) {
		return code
	}
	return http.StatusInternalServerError
}

func ErrorCode(code int) int {
	if IsErrorCode(code) {
		return code
	}
	return http.StatusInternalServerError
}

func IsStatusCode(code int) bool {
	return isInOpenRange(code, 99, 600)
}

func IsErrorCode(code int) bool {
	return isInOpenRange(code, 399, 600)
}

func IsClientErrorCode(code int) bool {
	return isInOpenRange(code, 399, 500)
}

func IsServerErrorCode(code int) bool {
	return isInOpenRange(code, 499, 600)
}

func isInOpenRange(subject int, min int, max int) bool {
	if subject < max && subject > min {
		return true
	}
	return false
}
